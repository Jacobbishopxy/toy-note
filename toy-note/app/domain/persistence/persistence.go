package persistence

import (
	"context"
	"fmt"
	"time"
	"toy-note/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "toy-note/entity"
)

var slog = logger.NewSugar("ToyNoteRepository")

// A ToyNoteRepository consists of two connections, one for sql and one for MongoDB
type ToyNoteRepository struct {
	db       *gorm.DB
	mongo    *mongo.Client
	mongoCtx context.Context
}

// ToyNoteRepository constructor
func NewToyNoteRepository(sqlUri, mongoUri string) (*ToyNoteRepository, error) {
	slog.Debug(fmt.Sprintf("Connecting to sql: %v", sqlUri))
	slog.Debug(fmt.Sprintf("Connecting to mongo: %v", mongoUri))

	db, err := gorm.Open(postgres.Open(sqlUri))
	if err != nil {
		return nil, err
	}

	slog.Debug("Connected to sql")

	mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	slog.Debug("Connected to mongo")

	return &ToyNoteRepository{
		db:       db,
		mongo:    mongo,
		mongoCtx: mongoCtx,
	}, nil
}

// Disconnect from sql & mongo
func (r *ToyNoteRepository) Disconnect() {
	sqlDb, err := r.db.DB()
	if err != nil {
		return
	}
	sqlDb.Close()

	r.mongo.Disconnect(r.mongoCtx)
}

// ToyNoteRepository must implement entity.NoteRepository
// var _ entity.ToyNoteRepo = &ToyNoteRepository{}
