package post

type GetPostsResponse struct {
	Posts []*Post
	Total int64 
	Page int64 
	Limit int64 
}

type CreatePostInput struct {
	Title string
	Content string
	Images [][]byte
}

type UpdatePostInput struct {
	UUID string
	Title string
	Content string
	Images [][]byte
	ImagesToDelete []string
}

type PostService interface {
	CreatePost(input *CreatePostInput) error
	GetPosts(page, limit int64) (*GetPostsResponse, error)
	GetPostByUUID(uuid string) (*Post, error)
	UpdatePost(input *UpdatePostInput) error
	DeletePost(uuid string) error
}

type PostRepository interface {
	Save(post *Post) error
	FindAll(page, limit int64) (*GetPostsResponse, error)
	FindByUUID(uuid string) (*Post, error)
	Update(post *Post) error
	Delete(uuid string) error
	DeleteImages(postUUID string, imagesUuids []string) error
}