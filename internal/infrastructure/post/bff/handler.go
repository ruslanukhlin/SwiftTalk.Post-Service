package bff

import (
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

func (h *Handler) GetPosts(c *fiber.Ctx) error {
	posts, err := h.postService.GetPosts(c.Context())
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(fiber.Map{
		"posts": posts,
	})
}

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

