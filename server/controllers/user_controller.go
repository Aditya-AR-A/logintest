package controllers

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"

	"github.com/aditya/logintest3/database"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



func RegisterUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Failed to parse user data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user data"})
	}

	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password cannot be empty"})
	}

	// Start a transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	defer tx.Rollback()

	// Check for existing user
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM login_credentials WHERE username = ? OR email = ?", user.Username, user.Email).Scan(&count)
	if err != nil {
		log.Printf("Database query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username or email already exists"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process password"})
	}

	// Insert new user
	_, err = tx.Exec("INSERT INTO login_credentials (username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Email, string(hashedPassword), time.Now(), time.Now())
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	log.Printf("User registered successfully: %s", user.Email)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}


func LoginUser(c *fiber.Ctx) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid login data"})
	}

	if loginData.Email == "" || loginData.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	var user User
	err := database.DB.QueryRow("SELECT id, username, email, password FROM login_credentials WHERE email = ?", loginData.Email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		log.Printf("Database query error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get session"})
	}
	sess.Set("user_id", user.ID)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
	}

	user.Password = "" // Remove password from response
	return c.JSON(fiber.Map{"message": "Login successful", "user": user})}



func LogoutUser(c *fiber.Ctx) error {
	sess, ok := c.Locals("session").(*session.Session)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get session"})
	}
	
	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to destroy session"})
	}
	return c.JSON(fiber.Map{"message": "Logout successful"})}



func GetAllUsers(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id, username, email, created_at, updated_at FROM login_credentials")
	if err != nil {
		log.Printf("Failed to query users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("Failed to scan user: %v", err)
			continue
		}
		users = append(users, user)
	}

	return c.JSON(users)
}


func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := database.DB.Exec("DELETE FROM login_credentials WHERE id = ?", id)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}



func ResetPassword(c *fiber.Ctx) error {
	email := c.FormValue("email")
	newPassword := c.FormValue("new_password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process new password"})
	}

	result, err := database.DB.Exec("UPDATE login_credentials SET password = ? WHERE email = ?", string(hashedPassword), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}
