package postgres

import (
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"gorm.io/gorm"
)

var _ domain.PostRepository = &PostgresMemoryRepository{}

type PostgresMemoryRepository struct {
	db *gorm.DB
}

func NewPostgresMemoryRepository(db *gorm.DB) *PostgresMemoryRepository {
	return &PostgresMemoryRepository{
		db: db,
	}
}

func (r *PostgresMemoryRepository) Save(post *domain.Post) error {
	var postDb Post
	
	postDb.UUID = post.UUID
	postDb.Title = post.Title.Value
	postDb.Content = post.Content.Value
	postDb.CreatedAt = post.CreatedAt
	postDb.UpdatedAt = post.UpdatedAt

	return r.db.Create(&postDb).Error
}

func (r *PostgresMemoryRepository) FindAll() ([]*domain.Post, error) {
	var posts []*Post
	if err := r.db.Find(&posts).Error; err != nil {
		return nil, err
	}

	domainPosts := make([]*domain.Post, len(posts))
	for i, post := range posts {
		domainPosts[i] = &domain.Post{
			UUID:      post.UUID,
			Title:     domain.Title{Value: post.Title},
			Content:   domain.Content{Value: post.Content},
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	return domainPosts, nil
}

func (r *PostgresMemoryRepository) FindByUUID(uuid string) (*domain.Post, error) {
	var post *Post
	if err := r.db.Where("uuid = ?", uuid).First(&post).Error; err != nil {
		return nil, err
	}

	return &domain.Post{
		UUID:      post.UUID,
		Title:     domain.Title{Value: post.Title},
		Content:   domain.Content{Value: post.Content},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (r *PostgresMemoryRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *PostgresMemoryRepository) Delete(uuid string) error {
	return r.db.Delete(&Post{}, uuid).Error
}