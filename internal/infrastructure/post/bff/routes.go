package bff

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	bffGroup := app.Group("post")

	bffGroup.Get("/", h.GetPosts)
	bffGroup.Post("/", h.CreatePost)
} 