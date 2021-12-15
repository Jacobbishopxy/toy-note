package adapters

import (
	"context"
	"fmt"
	"time"
	"toy-note/app/domain"
	"toy-note/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "toy-note/domain"
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

// ToyNoteRepository must implement domain.NoteRepository
// var _ domain.ToyNoteRepo = &ToyNoteRepository{}

// Create a new tag
func (r *ToyNoteRepository) SaveTag(tag *domain.Tag) error {
	return r.db.Debug().Create(tag).Error
}

// Get all tags
func (r *ToyNoteRepository) GetAllTags() ([]*domain.Tag, error) {
	var tags []*domain.Tag
	return tags, r.db.Debug().Find(&tags).Error
}

// Update tag by its id
func (r *ToyNoteRepository) UpdateTag(tag *domain.Tag) error {
	return r.db.Debug().Model(&tag).Updates(tag).Error
}

// Delete tag by its id
func (r *ToyNoteRepository) DeleteTag(id uint) error {
	return r.db.Debug().Delete(&domain.Tag{}, id).Error
}

// Create a new note
func (r *ToyNoteRepository) SavePost(post *domain.Post) error {
	return r.db.Debug().Create(post).Error
}
