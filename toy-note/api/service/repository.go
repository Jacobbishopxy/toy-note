package service

import (
	"io"
	"toy-note/api/entity"
)

// ToyNoteRepo
//
// Define an interface for a ToyNote repository
type ToyNoteRepo interface {
	// Get all tags
	GetTags() ([]entity.Tag, error)

	// Create/Update a tag
	// - If the tag Id is null, create a new tag
	// - If the tag Id is not null, update the existing tag
	SaveTag(tag entity.Tag) (entity.Tag, error)

	// Delete an existing tag
	DeleteTag(uint) error

	// Get posts by pagination
	GetPosts(entity.Pagination) ([]entity.Post, error)

	// Create/Update a post
	// - If the post Id is null, create a new post
	// - If the post Id is not null, update the existing post
	SavePost(entity.Post) (entity.Post, error)

	// Delete an existing post
	DeletePost(uint) error

	// Upload an affiliate
	UploadAffiliate(io.Reader, string) (string, error)

	// Download an affiliate
	DownloadAffiliate(uint) (entity.FileObject, error)

	// [admin] Get all unowned affiliates by pagination
	GetUnownedAffiliates(entity.Pagination) ([]entity.Affiliate, error)

	// [admin] Rebind a unowned affiliate to a post
	RebindAffiliate(uint, uint) error

	// [admin] Remove affiliates, which will remove affiliates files from mongo as well
	DeleteUnownedAffiliates([]uint) error

	// Search posts by tags
	SearchPostsByTags([]uint, entity.Pagination) ([]entity.Post, error)

	// TODO:
	// Search posts by title
}
