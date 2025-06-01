package application

import domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"

var _ domain.PostService = &PostApp{}

type PostApp struct {
	domain.PostRepository
}

func NewPostApp(postRepo domain.PostRepository) *PostApp {
	return &PostApp{
		PostRepository: postRepo,
	}
}

func (a *PostApp) CreatePost(post *domain.Post) error {
	return a.PostRepository.Save(post)
}

func (a *PostApp) GetPosts() ([]*domain.Post, error) {
	return a.PostRepository.FindAll()
}

func (a *PostApp) GetPostByUUID(uuid string) (*domain.Post, error) {
	return a.PostRepository.FindByUUID(uuid)
}

func (a *PostApp) UpdatePost(post *domain.Post) error {
	return a.PostRepository.Update(post)
}

func (a *PostApp) DeletePost(uuid string) error {
	return a.PostRepository.Delete(uuid)
}