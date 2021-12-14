package domain

type NoteRepository interface {
	SaveTag(*Tag) (*Tag, error)
	GetAllTags() ([]*Tag, error)
	UpdateTag(*Tag) (*Tag, error)
	DeleteTagByID(uint) error

	SavePost(*Post) (*Post, error)
	GetAllPosts() ([]*Post, error)
	UpdatePost(*Post) (*Post, error)
	DeletePostByID(uint) error
}
