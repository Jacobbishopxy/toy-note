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
	Host    string
	Port    int
	User    string
	Pass    string
	Db      string
	Sslmode string
}

// constructor
func NewPgRepository(logger *logger.ToyNoteLogger, conn PgConn) (PgRepository, error) {
	sqlUri := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s Timezone=Asia/Shanghai",
		conn.Host,
		conn.Port,
		conn.User,
		conn.Pass,
		conn.Db,
		conn.Sslmode,
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
	err := r.db.Exec("TRUNCATE TABLE posts, tags, affiliates RESTART IDENTITY CASCADE;").Error
	r.logger.Debug(fmt.Sprintf("TruncateAll: %v", err))
	return err
}

// This interface only denotes methods of pg repository should be implemented
type pgRepositoryInterface interface {
	// Get all tags at once
	GetTags() ([]entity.Tag, error)

	// Get a tag by id
	GetTag(uint) (entity.Tag, error)

	// Create a new tag, and return the created tag with id
	CreateTag(entity.Tag) (entity.Tag, error)

	// Update an existing tag, based on id
	UpdateTag(entity.Tag) (entity.Tag, error)

	// Delete an existing tag by id
	DeleteTag(uint) error

	// Get posts by pagination, ordered by created_at desc
	GetPosts(entity.Pagination) ([]entity.Post, error)

	// Get a post by id, including tags and affiliates
	GetPost(uint) (entity.Post, error)

	// Create a new post, and associate it with existing tags and affiliates
	CreatePost(entity.Post) (entity.Post, error)

	// Update an existing post, tags and affiliates are updated as well.
	// Using `Association` to deal with tags and affiliates, which means that
	// all the tags and affiliates should always be given in the request.
	// Any tags or affiliates was previously given and not given by now will be
	// unbounded from the post. They will not be deleted, so later if we need
	// them to be appeared in the post, we can still bind them to the post.
	UpdatePost(entity.Post) (entity.Post, error)

	// Delete an existing post, disassociate it with all tags and affiliates
	DeletePost(uint) error

	// Create/Update a new affiliate, notice that the affiliate don't need to be
	// associated to any post.
	// This method should not be exposed to the user.
	SaveAffiliate(entity.Affiliate) (entity.Affiliate, error)

	// Find an affiliate by id
	GetAffiliate(uint) (entity.Affiliate, error)

	// Find all unowned affiliates by ids
	GetUnownedAffiliatesByIds([]uint) ([]entity.Affiliate, error)

	// Find all unowned affiliates by pagination
	GetUnownedAffiliates(entity.Pagination) ([]entity.Affiliate, error)

	// Delete an unowned affiliate
	DeleteUnownedAffiliates([]uint) error

	// Find posts by tags
	GetPostsByTags([]uint, entity.Pagination) ([]entity.Post, error)

	// Find posts by title
	GetPostsByTitle(string, entity.Pagination) ([]entity.Post, error)

	// Find posts by time range
	GetPostsByTimeRange(entity.TimeSearch, entity.Pagination) ([]entity.Post, error)
}

var _ pgRepositoryInterface = (*PgRepository)(nil)

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
	result := r.db.Delete(entity.Tag{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("tag %d not found", id)
	}
	return result.Error
}

// ============================================================================
// Post
// ============================================================================

// calculate limit and offset for db query
func paginationToLimitOffset(pagination entity.Pagination) (int, int) {
	offset := (pagination.Page - 1) * pagination.Size
	return pagination.Size, offset
}

// private method
func (r *PgRepository) getPosts(ids []uint, pagination entity.Pagination) ([]entity.Post, error) {
	var posts []entity.Post
	// calc limit & offset
	limit, offset := paginationToLimitOffset(pagination)
	// preload all associations so that each post would be filled with tags and affiliates;
	// otherwise, the tags and affiliates would be empty
	que := r.db.
		Preload(clause.Associations).
		Limit(limit).
		Offset(offset)

	var err error

	if len(ids) == 0 {
		err = que.Find(&posts).Error
	} else {
		err = que.Where("id IN ?", ids).Find(&posts).Error
	}

	if err != nil {
		return posts, err
	}

	return posts, nil
}

func (r *PgRepository) GetPosts(pagination entity.Pagination) ([]entity.Post, error) {
	return r.getPosts(nil, pagination)
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
	// transaction here to make sure all the data modification is atomic
	r.db.Transaction(func(tx *gorm.DB) error {
		// update tags by replacing it
		if err := tx.Model(&post).Association("Tags").Replace(post.Tags); err != nil {
			return err
		}

		// update affiliates by replacing it
		if err := tx.Model(&post).Association("Affiliates").Replace(post.Affiliates); err != nil {
			return err
		}

		// update post
		if err := tx.Model(&post).Updates(post).Error; err != nil {
			return err
		}
		return nil
	})

	return post, nil
}

func (r *PgRepository) DeletePost(id uint) error {
	// transaction here to make sure all the data deletion is atomic
	return r.db.Transaction(func(tx *gorm.DB) error {

		var post entity.Post
		if err := tx.First(&post, id).Error; err != nil {
			return err
		}

		// do not delete data, but unbound from the post
		if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
			return err
		}

		// do not delete data, but unbound from the post
		if err := tx.Model(&post).Association("Affiliates").Clear(); err != nil {
			return err
		}

		// delete post
		if err := tx.Delete(&post).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *PgRepository) SaveAffiliate(affiliate entity.Affiliate) (entity.Affiliate, error) {
	if err := r.db.Save(&affiliate).Error; err != nil {
		return entity.Affiliate{}, err
	}
	return affiliate, nil
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

type PostsTags struct {
	PostID uint
}

const searchByTagQuery = `
SELECT
	post_id
FROM
	posts_tags
WHERE
	tag_id IN ?
GROUP BY
	post_id
HAVING
	count(distinct tag_id) = ?
`

func (r *PgRepository) GetPostsByTags(
	tagIds []uint,
	pagination entity.Pagination,
) ([]entity.Post, error) {
	var postIds []uint

	err := r.db.
		Raw(searchByTagQuery, tagIds, len(tagIds)).
		Scan(&postIds).
		Error
	if err != nil {
		return nil, err
	}

	return r.getPosts(postIds, pagination)
}

func (r *PgRepository) GetPostsByTitle(
	title string,
	pagination entity.Pagination,
) ([]entity.Post, error) {
	var postIds []uint

	err := r.db.
		Raw("SELECT id FROM posts WHERE title LIKE ?", "%"+title+"%").
		Scan(&postIds).
		Error
	if err != nil {
		return nil, err
	}

	return r.getPosts(postIds, pagination)
}

func (r *PgRepository) GetPostsByTimeRange(
	timeSearch entity.TimeSearch,
	pagination entity.Pagination,
) ([]entity.Post, error) {
	var postIds []uint

	var d string
	switch timeSearch.Type {
	case 0:
		d = "date"
	case 1:
		d = "created_at"
	case 2:
		d = "updated_at"
	default:
		d = "date"
	}

	err := r.db.
		Raw(
			"SELECT id FROM posts WHERE ? BETWEEN ? AND ?",
			d,
			timeSearch.Start,
			timeSearch.End,
		).
		Scan(&postIds).
		Error
	if err != nil {
		return nil, err
	}

	return r.getPosts(postIds, pagination)
}
