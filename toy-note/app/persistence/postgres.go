package persistence

import (
	"fmt"
	"toy-note/app/entity"
	"toy-note/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Postgresql Repository
type PgRepository struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

// constructor
func NewPgRepository(logger *logger.ToyNoteLogger, sqlUri string) (PgRepository, error) {
	slog := logger.NewSugar("PgRepository")

	slog.Debug(fmt.Sprintf("Connecting to sql: %v", sqlUri))

	db, err := gorm.Open(postgres.Open(sqlUri))
	if err != nil {
		return PgRepository{}, err
	}

	slog.Debug("Connected to sql")

	return PgRepository{
		logger: slog,
		db:     db,
	}, nil
}

// Auto Migrate. Create tables if not exists
func (r *PgRepository) AutoMigrate() error {
	err := r.db.AutoMigrate(&entity.Tag{}, &entity.Affiliate{}, &entity.Post{})
	r.logger.Debug(fmt.Sprintf("AutoMigrate: %v", err))
	return err
}

// This interface only denotes methods of pg repository should be implemented
type pgRepositoryInterface interface {
	// Get all tags at once
	GetTags() ([]entity.Tag, error)

	// Get a tag by id
	GetTag(uint) (entity.Tag, error)

	// Create a new tag
	CreateTag(tag entity.Tag) error

	// Update an existing tag
	UpdateTag(tag entity.Tag) error

	// Delete an existing tag
	DeleteTag(uint) error

	// Get posts by pagination
	GetPosts(entity.Pagination) ([]entity.Post, error)

	// Get a post by id, including tags and affiliates
	GetPost(uint) (entity.Post, error)

	// Create a new post, and associate it with existing tags and affiliates
	CreatePost(post entity.Post) error

	// Update an existing post, tags and affiliates can be updated as well
	UpdatePost(post entity.Post) error

	// Delete an existing post, disassociate it with all tags and delete affiliates
	DeletePost(uint) error
}

var _ pgRepositoryInterface = &PgRepository{}

// ============================================================================
// Tag
// ============================================================================

func (r *PgRepository) GetTags() ([]entity.Tag, error) {
	var tags []entity.Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *PgRepository) GetTag(id uint) (entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return tag, err
	}

	return tag, nil
}

func (r *PgRepository) CreateTag(tag entity.Tag) error {
	return r.db.Create(&tag).Error
}

func (r *PgRepository) UpdateTag(tag entity.Tag) error {
	return r.db.Save(&tag).Error
}

func (r *PgRepository) DeleteTag(id uint) error {
	return r.db.Delete(entity.Tag{}, id).Error
}

// ============================================================================
// Post
// ============================================================================

func (r *PgRepository) GetPosts(pagination entity.Pagination) ([]entity.Post, error) {
	var posts []entity.Post

	offset := (pagination.Page - 1) * pagination.Size
	err := r.db.
		Preload(clause.Associations).
		Limit(pagination.Size).
		Offset(offset).
		Association("Tags").
		Find(&posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PgRepository) GetPost(id uint) (entity.Post, error) {
	var post entity.Post
	err := r.db.Preload(clause.Associations).First(&post, id).Error
	if err != nil {
		return post, err
	}

	return post, nil
}

func (r *PgRepository) CreatePost(post entity.Post) error {
	return r.db.Save(&post).Error
}

func (r *PgRepository) UpdatePost(post entity.Post) error {
	return r.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&post).
		Error
}

func (r *PgRepository) DeletePost(id uint) error {
	return r.db.Select("Affiliates").Delete(entity.Post{}, id).Error
}
