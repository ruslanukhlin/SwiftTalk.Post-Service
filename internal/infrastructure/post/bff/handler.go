package bff

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	postService *PostService
}

func NewHandler(postService *PostService) *Handler {
	return &Handler{
		postService: postService,
	}
}

// GetPost godoc
// @Summary Получить пост по ID
// @Description Получить детальную информацию о посте по его идентификатору
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "ID поста"
// @Success 200 {object} fiber.Map "Успешный ответ с информацией о посте"
// @Failure 400 {object} ErrorResponse "Ошибка запроса"
// @Failure 404 {object} ErrorResponse "Пост не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /post/{id} [get]
func (h *Handler) GetPost(c *fiber.Ctx) error {
	postId := c.Params("uuid")

	post, err := h.postService.GetPost(c, postId)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"post": post,
	})
}

// GetPosts godoc
// @Summary Получить список постов
// @Description Получить список всех постов
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} GetPostsResponse "Успешный ответ с списком постов"
// @Failure 400 {object} ErrorResponse "Ошибка запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /post [get]
func (h *Handler) GetPosts(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}
	
	posts, err := h.postService.GetPosts(c, pageInt, limitInt)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"posts": posts.Posts,
		"total": posts.Total,
		"page": posts.Page,
		"limit": posts.Limit,
	})
}

// CreatePost godoc
// @Summary Создать новый пост
// @Description Создать новый пост с заголовком, содержанием и опциональными изображениями
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Заголовок поста"
// @Param content formData string true "Содержание поста"
// @Param images formData file false "Изображения (множественная загрузка)" collectionFormat="multi"
// @Success 200 {object} CreatePostResponse "Успешное создание поста"
// @Failure 400 {object} ErrorResponse "Ошибка в параметрах запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /post [post]
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	if title == "" || content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	images := form.File["images"]

	err = h.postService.CreatePost(c, title, content, images)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Пост успешно создан",
	})
}

// UpdatePost godoc
// @Summary Обновить пост по ID
// @Description Обновить пост по его идентификатору
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param uuid path string true "ID поста"
// @Param title formData string true "Заголовок поста"
// @Param content formData string true "Содержание поста"
// @Param images formData file false "Изображения (множественная загрузка)" collectionFormat="multi"
// @Param deleted_images formData string false "Удаленные изображения" collectionFormat="multi"
// @Success 200 {object} UpdatePostResponse "Успешное обновление поста"
// @Failure 400 {object} ErrorResponse "Ошибка в параметрах запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /post/{id} [put]
func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	postId := c.Params("uuid")
	title := c.FormValue("title")
	content := c.FormValue("content")
	deletedImages := c.FormValue("deleted_images")

	var imageUUIDs []string
	if deletedImages != "" {
		if err := json.Unmarshal([]byte(deletedImages), &imageUUIDs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Неверный формат данных для deleted_images. Ожидается JSON массив строк.",
			})
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	images := form.File["images"]

	if title == "" || content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	err = h.postService.UpdatePost(c, postId, title, content, images, imageUUIDs)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Пост успешно обновлен",
	})
}	

// DeletePost godoc
// @Summary Удалить пост по ID
// @Description Удалить пост по его идентификатору
// @Tags posts
// @Accept json
// @Produce json
// @Param id path string true "ID поста"
// @Success 200 {object} fiber.Map "Успешное удаление поста"
// @Failure 400 {object} ErrorResponse "Ошибка запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /post/{id} [delete]
func (h *Handler) DeletePost(c *fiber.Ctx) error {
	postId := c.Params("uuid")

	err := h.postService.DeletePost(c, postId)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Пост успешно удален",
	})
}