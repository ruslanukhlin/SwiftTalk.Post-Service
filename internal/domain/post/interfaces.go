package domain

type PostService interface {
	CreatePost(post *Post) error
	GetPosts() ([]*Post, error)
	GetPostByUUID(uuid string) (*Post, error)
	UpdatePost(post *Post) error
	DeletePost(uuid string) error
}

type PostRepository interface {
	Save(post *Post) error
	FindAll() ([]*Post, error)
	FindByUUID(uuid string) (*Post, error)
	Update(post *Post) error
	Delete(uuid string) error
}