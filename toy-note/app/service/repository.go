package service

import (
	"toy-note/app/entity"
)

// ToyNoteRepo
//
// Define an interface for a ToyNote repository
type ToyNoteRepo interface {
	// Get all tags
	GetTags() ([]entity.Tag, error)

	// Create a new tag
	CreateTag(tag entity.Tag) (*entity.Tag, error)

	// Update an existing tag
	UpdateTag(tag entity.Tag) (*entity.Tag, error)

	// Delete an existing tag
	DeleteTag(tag entity.Tag) error

	// Get notes by pagination
	GetByPagination(entity.Pagination) []entity.Post

	// Create a new note
	CreateNote(title, context string) (*entity.Post, error)

	// Update an existing note
	UpdateNote(entity.Post) (*entity.Post, error)

	// Delete an existing note
	DeleteNote(entity.Post) error

	// Download an affiliate
	DownloadAffiliate(string) ([]byte, error)
}
