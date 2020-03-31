package main

func initializeRoutes() {
	router.GET("/status", statusResponseHandler)
	router.POST("/start", startResponseHander)
}
