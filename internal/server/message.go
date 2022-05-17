package server

import (
	"net/http"

	"github.com/DATATRONiQ/go-sparkplug-primary/internal/store"
	"github.com/gin-gonic/gin"
)

func indexMessages(ctx *gin.Context) {
	messages := store.FetchMessages()
	ctx.JSON(http.StatusOK, gin.H{
		"data": messages,
	})
}
