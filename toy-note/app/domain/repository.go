package domain

type Pagination struct {
	Page  int
	Limit int
}

// ToyNoteRepo
//
// Define an interface for a ToyNote repository
type ToyNoteRepo interface {
	GetByPagination(Pagination) []Note
	New(title, context string) (Note, error)
	Update(Note) (Note, error)
}
