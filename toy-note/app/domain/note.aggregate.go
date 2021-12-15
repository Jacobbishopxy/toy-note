package domain

type Note struct {
	post       Post
	tags       []Tag
	affiliates []Affiliate
}

func NewNote(title, context string) Note {
	return Note{
		post: Post{
			Title:   title,
			Content: context,
		},
		tags:       []Tag{},
		affiliates: []Affiliate{},
	}
}

func (n *Note) GetPost() Post {
	return n.post
}

func (n *Note) GetTags() []Tag {
	return n.tags
}

func (n *Note) GetAffiliates() []Affiliate {
	return n.affiliates
}

// TODO: the rest of the methods are required by business logic
