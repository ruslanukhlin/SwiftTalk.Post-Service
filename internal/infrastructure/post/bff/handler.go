package bff

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	postService *PostService
}

// PostsResponse представляет ответ со списком постов
type PostsResponse struct {
	Posts []Post `json:"posts"`
}

// CreatePostResponse представляет ответ на создание поста
type CreatePostResponse struct {
	Message string `json:"message"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(postService *PostService) *Handler {
	return &Handler{
		postService: postService,
	}
}

// GetPosts godoc
// @Summary Получить список постов
// @Description Получить список всех постов
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} PostsResponse "Успешный ответ с списком постов"
// @Failure 400 {object} ErrorResponse "Ошибка запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /post [get]
func (h *Handler) GetPosts(c *fiber.Ctx) error {
	posts, err := h.postService.GetPosts(c.Context())
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"posts": posts,
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

	err = h.postService.CreatePost(c.Context(), title, content, images)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Пост успешно создан",
	})
}

