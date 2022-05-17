package server

import (
	"net/http"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	"github.com/gin-gonic/gin"
)

func setRouter(sm *store.StoreManager) *gin.Engine {
	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()

	// Create API route group
	api := router.Group("/api")
	{
		// Add /hello GET route to router and define route handler function
		api.GET("/hello", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"msg": "world"})
		})
	}

	api.GET("/messages", indexMessages)
	api.GET("/groups", func(ctx *gin.Context) {
		groups := sm.Fetch()
		ctx.JSON(http.StatusOK, gin.H{
			"data": groups,
		})
	})

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	return router
}
