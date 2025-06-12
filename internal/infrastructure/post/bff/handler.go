package bff

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": st.Message(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Внутренняя ошибка сервера",
				})
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"posts": posts,
	})
}

func (h *Handler) CreatePost(c *fiber.Ctx) error {
	payload := new(CreatePostPayload)
	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат данных",
		})
	}

	err := h.postService.CreatePost(c.Context(), payload)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": st.Message(),
				})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Внутренняя ошибка сервера",
				})
			}
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Пост успешно создан",
	})
}

