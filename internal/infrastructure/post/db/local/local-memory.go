package db

import (
	"errors"
	"sync"

	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
)

var (
	ErrPostNotFound = errors.New("post not found")
)

var _ domain.PostRepository = &LocalMemoryPostRepository{}

type LocalMemoryPostRepository struct {
	posts map[string]*domain.Post
	mu sync.RWMutex
}

func NewLocalMemoryPostRepository() *LocalMemoryPostRepository {
	return &LocalMemoryPostRepository{
		posts: make(map[string]*domain.Post),
	}
}

func (r *LocalMemoryPostRepository) Save(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.UUID] = post

	return nil
}

func (r *LocalMemoryPostRepository) FindAll() ([]*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	posts := make([]*domain.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *LocalMemoryPostRepository) FindByUUID(uuid string) (*domain.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, ok := r.posts[uuid]
	if !ok {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (r *LocalMemoryPostRepository) Update(post *domain.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.UUID] = post

	return nil
}

func (r *LocalMemoryPostRepository) Delete(uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.posts, uuid)

	return nil
}