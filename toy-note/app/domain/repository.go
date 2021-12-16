package domain

type Pagination struct {
	Page int
	Size int
}

// ToyNoteRepo
//
// Define an interface for a ToyNote repository
type ToyNoteRepo interface {
	// Get all tags
	GetTags() ([]Tag, error)

	// Create a new tag
	CreateTag(tag Tag) (*Tag, error)

	// Update an existing tag
	UpdateTag(tag Tag) (*Tag, error)

	// Delete an existing tag
	DeleteTag(tag Tag) error

	// Get notes by pagination
	GetByPagination(Pagination) []Note

	// Create a new note
	CreateNote(title, context string) (*Note, error)

	// Update an existing note
	UpdateNote(Note) (*Note, error)

	// Delete an existing note
	DeleteNote(Note) error

	// Download an affiliate
	DownloadAffiliate(string) ([]byte, error)
}
