package persistence

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"toy-note/api/entity"
	"toy-note/errors"
	"toy-note/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const DatabaseName = "toy-note"
const CollectionName = "fs.files"

// A MongoRepository consists of two connections, one for sql and one for MongoDB
type MongoRepository struct {
	logger *zap.SugaredLogger
	db     *mongo.Database
}

type MongoConn struct {
	host string
	port int
	user string
	pass string
}

// constructor
func NewMongoRepository(ctx context.Context, logger *logger.ToyNoteLogger, conn MongoConn) (MongoRepository, error) {
	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%d", conn.user, conn.pass, conn.host, conn.port)

	slog := logger.NewSugar("MongoRepository")
	slog.Debug(fmt.Sprintf("Connecting to mongo: %v", mongoUri))

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return MongoRepository{}, errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to connect to MongoDB")
	}
	db := client.Database(DatabaseName)

	slog.Debug("Connected to mongo")

	return MongoRepository{
		logger: slog,
		db:     db,
	}, nil
}

// Make sure the MongoRepository implements the Repository interface
type mongoRepositoryInterface interface {
	UploadFile(reader io.Reader, filename string) (string, error)
	DownloadFile(filename, id string) (entity.FileObject, error)
	DeleteFiles(ids []string) error
}

var _ mongoRepositoryInterface = &MongoRepository{}

// Upload file to MongoDB
// The result string is the object id from the MongoDB, which is supposed to be stored in PG.
func (r *MongoRepository) UploadFile(reader io.Reader, filename string) (string, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to read file")
	}

	bucket, err := gridfs.NewBucket(r.db)
	if err != nil {
		return "", errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to create bucket")
	}

	uploadStream, err := bucket.OpenUploadStream(filename)
	if err != nil {
		return "", errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to open upload stream")
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		return "", errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to write data to upload stream")
	}

	return uploadStream.FileID.(primitive.ObjectID).Hex(), nil
}

// Download file from MongoDB, according to the id
func (r *MongoRepository) DownloadFile(filename, id string) (entity.FileObject, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return entity.FileObject{}, errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to convert id to ObjectID")
	}

	bucket, err := gridfs.NewBucket(r.db)
	if err != nil {
		return entity.FileObject{}, errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to create bucket")
	}

	var buf bytes.Buffer
	size, err := bucket.DownloadToStream(oid, &buf)
	if err != nil {
		return entity.FileObject{}, errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to download stream")
	}

	r.logger.Debug(fmt.Sprintf("File download completed, size: %v", size))

	return entity.FileObject{
		Filename: filename,
		Content:  buf.Bytes(),
		Size:     size,
	}, nil
}

// Delete files from MongoDB, according to the ids
func (r *MongoRepository) DeleteFiles(ids []string) error {
	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to convert id to ObjectID")
		}

		oids = append(oids, oid)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := r.db.Collection(CollectionName).DeleteMany(ctx, bson.M{"_id": bson.M{"$in": oids}}); err != nil {
		return errors.WrapErrorf(err, errors.ErrorCodeUnknown, "Failed to delete files")
	}

	return nil
}
