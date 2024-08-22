// Package router provides routing and middleware setup for the application.
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/aditya/logintest3/server/controllers"
)

// store is a global session store for managing user sessions.
var store *session.Store

// SetupRouter configures and returns a Fiber application with all routes and middleware set up.
func SetupRouter() *fiber.App {
	app := fiber.New()

	// setupMiddleware configures all middleware for the application.
	setupMiddleware(app)

	// setupRoutes defines all the routes for the application.
	setupRoutes(app)

	return app
}

// setupMiddleware configures CORS, static file serving, and session handling middleware.
func setupMiddleware(app *fiber.App) {
	// CORS middleware configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	// Serve static files
	app.Static("/", "./media")

	// Initialize session store
	store = session.New()

	// Session middleware
	app.Use(func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get session"})
		}
		c.Locals("session", sess)
		return c.Next()
	})
}

// setupRoutes defines all the routes for the application.
func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./media/main.html")
	})
	app.Post("/register", controllers.RegisterUser)
	app.Post("/login", controllers.LoginUser)
	app.Post("/logout", controllers.LogoutUser)
	app.Get("/users", controllers.GetAllUsers)
	app.Delete("/users/:id", controllers.DeleteUser)
	app.Get("/users-page", func(c *fiber.Ctx) error {
		return c.SendFile("./media/users.html")
	})
	app.Post("/reset-password", controllers.ResetPassword)
	app.Get("/reset-password-page", func(c *fiber.Ctx) error {
		return c.SendFile("./media/reset_password.html")
	})
}