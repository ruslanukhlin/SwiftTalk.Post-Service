package bff

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/post", h.GetPosts)
	app.Post("/post", h.CreatePost)
	app.Get("/post/:uuid", h.GetPost)
	app.Delete("/post/:uuid", h.DeletePost)
	app.Patch("/post/:uuid", h.UpdatePost)
} 