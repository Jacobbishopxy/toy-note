package persistence

import (
	"fmt"
	"toy-note/api/entity"
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

type PgConn struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	sslmode  string
}

// constructor
func NewPgRepository(logger *logger.ToyNoteLogger, conn PgConn) (PgRepository, error) {
	sqlUri := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s Timezone=Asia/Shanghai",
		conn.host,
		conn.port,
		conn.user,
		conn.password,
		conn.dbname,
		conn.sslmode,
	)

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

func (r *PgRepository) TruncateAll() error {
	return r.db.Exec("TRUNCATE TABLE posts, tags, affiliates RESTART IDENTITY CASCADE;").Error
}

// This interface only denotes methods of pg repository should be implemented
type pgRepositoryInterface interface {
	// Get all tags at once
	GetTags() ([]entity.Tag, error)

	// Get a tag by id
	GetTag(uint) (entity.Tag, error)

	// Create a new tag
	CreateTag(entity.Tag) (entity.Tag, error)

	// Update an existing tag
	UpdateTag(entity.Tag) (entity.Tag, error)

	// Delete an existing tag
	DeleteTag(uint) error

	// Get posts by pagination
	GetPosts(entity.Pagination) ([]entity.Post, error)

	// Get a post by id, including tags and affiliates
	GetPost(uint) (entity.Post, error)

	// Create a new post, and associate it with existing tags and affiliates
	CreatePost(entity.Post) (entity.Post, error)

	// Update an existing post, tags and affiliates are updated as well
	// Using `Association` to deal with tags and affiliates, which means that
	// all the tags and affiliates should always be given in the request.
	// Any tags or affiliates was previously given and not given by now will be
	// unbounded from the post. It will not be deleted, so later when we need them,
	// we can still bind them to the post.
	UpdatePost(entity.Post) (entity.Post, error)

	// Delete an existing post, disassociate it with all tags and delete affiliates
	DeletePost(uint) error

	// Find an affiliate by id
	GetAffiliate(uint) (entity.Affiliate, error)

	// Find all unowned affiliates by ids
	GetUnownedAffiliatesByIds([]uint) ([]entity.Affiliate, error)

	// Find all unowned affiliates by pagination
	GetUnownedAffiliates(entity.Pagination) ([]entity.Affiliate, error)

	// Delete an unowned affiliate
	DeleteUnownedAffiliates([]uint) error
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

func (r *PgRepository) CreateTag(tag entity.Tag) (entity.Tag, error) {
	if err := r.db.Create(&tag).Error; err != nil {
		return entity.Tag{}, err
	}
	return tag, nil
}

func (r *PgRepository) UpdateTag(tag entity.Tag) (entity.Tag, error) {
	if err := r.db.Updates(&tag).Error; err != nil {
		return entity.Tag{}, err
	}
	return tag, nil
}

func (r *PgRepository) DeleteTag(id uint) error {
	return r.db.Delete(entity.Tag{}, id).Error
}

// ============================================================================
// Post
// ============================================================================

func (r *PgRepository) GetPosts(pagination entity.Pagination) ([]entity.Post, error) {
	var posts []entity.Post
	// calc offset for db query
	offset := (pagination.Page - 1) * pagination.Size
	// preload all associations so that each post would be filled with tags and affiliates;
	// otherwise, the tags and affiliates would be empty
	err := r.db.
		Preload(clause.Associations).
		Limit(pagination.Size).
		Offset(offset).
		Find(&posts).
		Error

	if err != nil {
		return posts, err
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

func (r *PgRepository) CreatePost(post entity.Post) (entity.Post, error) {
	if err := r.db.Save(&post).Error; err != nil {
		return entity.Post{}, err
	}
	return post, nil
}

func (r *PgRepository) UpdatePost(post entity.Post) (entity.Post, error) {

	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&post).Association("Tags").Replace(post.Tags); err != nil {
			return err
		}

		if err := tx.Model(&post).Association("Affiliates").Replace(post.Affiliates); err != nil {
			return err
		}

		if err := tx.Model(&post).Updates(post).Error; err != nil {
			return err
		}
		return nil
	})

	return post, nil
}

func (r *PgRepository) DeletePost(id uint) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		var post entity.Post
		if err := tx.First(&post, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
			return err
		}

		if err := tx.Model(&post).Association("Affiliates").Clear(); err != nil {
			return err
		}

		if err := tx.Delete(&post).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PgRepository) GetAffiliate(id uint) (entity.Affiliate, error) {
	var affiliate entity.Affiliate
	if err := r.db.First(&affiliate, id).Error; err != nil {
		return affiliate, err
	}

	return affiliate, nil
}

func (r *PgRepository) GetUnownedAffiliatesByIds(ids []uint) ([]entity.Affiliate, error) {
	var affiliates []entity.Affiliate

	err := r.db.
		Where("post_refer IS NULL").
		Find(&affiliates, ids).
		Error

	if err != nil {
		return affiliates, err
	}

	return affiliates, nil
}

func (r *PgRepository) GetUnownedAffiliates(pagination entity.Pagination) ([]entity.Affiliate, error) {
	var affiliates []entity.Affiliate
	offset := (pagination.Page - 1) * pagination.Size

	err := r.db.
		Limit(pagination.Size).
		Offset(offset).
		Where("post_refer IS NULL").
		Find(&affiliates).
		Error

	if err != nil {
		return nil, err
	}

	return affiliates, nil
}

func (r *PgRepository) DeleteUnownedAffiliates(ids []uint) error {
	return r.db.Where("post_refer IS NULL").Delete(entity.Affiliate{}, ids).Error
}
