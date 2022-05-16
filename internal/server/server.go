package server

func Start() {
	router := setRouter()

	// Start listening and serving requests
	router.Run(":8080")
}
