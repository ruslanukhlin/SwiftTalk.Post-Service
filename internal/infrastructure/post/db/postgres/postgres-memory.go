package postgres

import (
	"errors"

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

	for _, imageUrl := range post.Images {
		postDb.Images = append(postDb.Images, Image{
			UUID:     imageUrl.UUID,
			URL:      imageUrl.URL,
			PostUUID: post.UUID,
		})
	}

	return r.db.Create(&postDb).Error
}

func (r *PostgresMemoryRepository) FindAll(page, limit int64) (*domain.GetPostsResponse, error) {
	var posts []*Post
	if err := r.db.Preload("Images").Limit(int(limit)).Offset(int((page - 1) * limit)).Find(&posts).Error; err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.Model(&Post{}).Count(&total).Error; err != nil {
		return nil, err
	}

	domainPosts := make([]*domain.Post, len(posts))
	for i, post := range posts {
		images := getImages(post.Images)
		domainPosts[i] = &domain.Post{
			UUID:      post.UUID,
			Title:     domain.Title{Value: post.Title},
			Content:   domain.Content{Value: post.Content},
			Images:    images,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
	}

	return &domain.GetPostsResponse{
		Posts: domainPosts,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (r *PostgresMemoryRepository) FindByUUID(uuid string) (*domain.Post, error) {
	var post Post
	if err := r.db.Preload("Images").Where("uuid = ?", uuid).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPostNotFound
		}
		return nil, err
	}

	images := getImages(post.Images)

	return &domain.Post{
		UUID:      post.UUID,
		Title:     domain.Title{Value: post.Title},
		Content:   domain.Content{Value: post.Content},
		Images:    images,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}, nil
}

func (r *PostgresMemoryRepository) Update(post *domain.Post) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var postDb Post
		if err := tx.Where("uuid = ?", post.UUID).First(&postDb).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrPostNotFound
			}
			return err
		}

		// Обновляем основные поля поста
		updates := map[string]interface{}{
			"title":   post.Title.Value,
			"content": post.Content.Value,
		}
		if err := tx.Model(&postDb).Updates(updates).Error; err != nil {
			return err
		}

		// Если есть новые изображения, обрабатываем их
		if post.Images != nil {
			// Создаем map существующих изображений для быстрого поиска
			existingImages := make(map[string]bool)
			for _, img := range postDb.Images {
				existingImages[img.UUID] = true
			}

			// Обновляем или добавляем новые изображения
			for _, newImage := range post.Images {
				if newImage == nil {
					continue
				}

				if !existingImages[newImage.UUID] {
					// Добавляем новое изображение
					if err := tx.Create(&Image{
						UUID:     newImage.UUID,
						URL:      newImage.URL,
						PostUUID: post.UUID,
					}).Error; err != nil {
						return err
					}
				}
			}

			// Удаляем изображения, которых нет в новом списке
			newImageMap := make(map[string]bool)
			for _, img := range post.Images {
				if img != nil {
					newImageMap[img.UUID] = true
				}
			}

			for _, oldImage := range postDb.Images {
				if !newImageMap[oldImage.UUID] {
					if err := tx.Delete(&Image{UUID: oldImage.UUID}).Error; err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
}

func (r *PostgresMemoryRepository) Delete(uuid string) error {
	err := r.db.Delete(&Post{UUID: uuid}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrPostNotFound
		}
		return err
	}

	return nil
}

func (r *PostgresMemoryRepository) DeleteImages(postUUID string, imagesUuids []string) error {
	return r.db.Where("post_uuid = ? AND uuid IN (?)", postUUID, imagesUuids).Delete(&Image{}).Error
}

func getImages(images []Image) []*domain.Image {
	imagesDomain := make([]*domain.Image, len(images))
	for i, image := range images {
		imagesDomain[i] = &domain.Image{
			UUID: image.UUID,
			URL:  image.URL,
		}
	}
	return imagesDomain
}
