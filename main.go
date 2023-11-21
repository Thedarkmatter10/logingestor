package main

import "logingestor/routes"

func main() {

	// Set up the router
	router := routes.SetUpRouter()

	// Start the HTTP server on port 3000
	router.Run(":3000")
}
