package bff

import (
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleGRPCError(c *fiber.Ctx, err error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": st.Message(),
			})
		case codes.NotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": st.Message(),
			})
		case codes.AlreadyExists:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
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