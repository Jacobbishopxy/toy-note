package domain

// ToyNoteRepo
//
// Define an interface for a ToyNote repository
type ToyNoteRepo interface {
	// Create a new tag
	SaveTag(*Tag) error

	// Get all tags
	GetAllTags() ([]*Tag, error)

	// Update tag by its id
	UpdateTag(*Tag) error

	// Delete tag by its id
	DeleteTag(uint) error

	// Create a new post.
	// Existing tags can be bound to the post.
	// Newly affiliate files can be saved as well.
	SavePost(*Post) error

	// Get all posts by pagination.
	// With tags and affiliate files.
	GetAllPosts(Pagination) ([]*Post, error)

	// Update a post by its id.
	// Existing tags can be bound to the post.
	// Newly affiliate files can be saved as well.
	UpdatePost(*Post) error

	// Delete a post by its id.
	DeletePost(uint) error

	// Bind existing tags to a post.
	BindTagsToPost([]uint, uint) error

	// Unbind tags from a post.
	UnbindTagsFromPost([]uint, uint) error

	// Upload a new affiliate file to a post.
	UploadAffiliateToPost(*PostAffiliate) (*PostAffiliate, error)

	// Download an affiliate file from a post.
	DownloadAffiliateFromPost(*PostAffiliate) error

	// Delete an affiliate file from a post.
	DeleteAffiliateFromPost(*PostAffiliate) error
}
