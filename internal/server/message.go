package server

import (
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	"github.com/gofiber/fiber/v2"
)

func indexMessages(ctx *fiber.Ctx) error {
	messages := store.Fetch()
	ctx.JSON(messages)
	return nil
}
