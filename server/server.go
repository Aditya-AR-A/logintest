// Package main is the entry point for the login system application.
package main

import (
	"log"

	"github.com/aditya/logintest3/database"
	"github.com/aditya/logintest3/server/router"
)

// main is the entry point of the application. It sets up the database connection,
// initializes the router, and starts the server.
func main() {
	// Connect to database
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// Setup and run the server
	app := router.SetupRouter()
	log.Fatal(app.Listen(":3000"))
}