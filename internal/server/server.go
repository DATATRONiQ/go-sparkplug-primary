package server

import "github.com/DATATRONiQ/go-sparkplug-primary/internal/store"

func Start(sm *store.StoreManager) {
	router := setRouter(sm)

	// Start listening and serving requests
	router.Listen(":8080")
}
