package server

import (
	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func setRouter(sm *store.StoreManager) *fiber.App {
	app := fiber.New()

	app.Static("/", "../../assets/build", fiber.Static{})

	// Create API route group
	apiGroup := app.Group("/api")

	apiGroup.Get("/messages", indexMessages)
	apiGroup.Get("/groups", func(ctx *fiber.Ctx) error {
		groups := sm.FetchFull()
		ctx.JSON(groups)
		return nil
	})

	apiGroup.Get("/groups/stream", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "text/event-stream")
		ctx.Set("Cache-Control", "no-cache,no-transform")
		ctx.Set("Connection", "keep-alive")
		ctx.Set("Access-Control-Allow-Origin", "*")
		ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(sm.GroupsSSEHandler.Subscribe))
		return nil
	})

	app.Static("*", "../../assets/build/index.html")

	return app
}
