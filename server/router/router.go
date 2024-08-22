package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/aditya/logintest3/server/controllers"
)

var store *session.Store

func SetupRouter() *fiber.App {
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

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

	// Routes
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

	return app
}
