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

type PgRepository struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

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

func (r *PgRepository) AutoMigrate() error {
	err := r.db.AutoMigrate(&entity.Tag{}, &entity.Affiliate{}, &entity.Post{})
	r.logger.Debug(fmt.Sprintf("AutoMigrate: %v", err))
	return err
}

type pgRepositoryInterface interface {
	GetTags() ([]entity.Tag, error)
	GetTag(uint) (entity.Tag, error)
	CreateTag(tag entity.Tag) error
	UpdateTag(tag entity.Tag) error
	DeleteTag(uint) error

	GetPosts(entity.Pagination) ([]entity.Post, error)
	GetPost(uint) (entity.Post, error)
	CreatePost(post entity.Post) error
	UpdatePost(post entity.Post) error
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
