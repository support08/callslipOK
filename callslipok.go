package main

import (
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve secret key from environment variable
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("SECRET_KEY not set in environment variables")
	}
	// Hardcode your values directly in the code
	branchID := "22345"
	apiKey := "SLIPOK6H4K8YP"
	// fileLocation := "C:/Support08 Work/NOTE Cgcloud app deployODE/sql tech/export/154210.jpg"
	app := fiber.New()

	// Middleware to check for Authorization header
	authMiddleware := func(c *fiber.Ctx) error {
		authHeader := c.Get("x-authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing Authorization header")
		}

		if authHeader != secretKey {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or missing API key")
		}

		return c.Next()
	}
	app.Post("/call-api", authMiddleware, func(c *fiber.Ctx) error {
		var requestBody struct {
			FilePath string `json:"file_path"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		// Ensure the FilePath field is not empty
		if requestBody.FilePath == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Empty file path")
		}
		client := resty.New()
		resp, err := client.R().
			SetHeader("x-authorization", apiKey).
			SetFile("files", requestBody.FilePath).
			SetFormData(map[string]string{
				"log": "true",
			}).
			Post("https://api.slipok.com/api/line/apikey/" + branchID)

		if err != nil {
			log.Fatalf("Error calling API: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error calling API")
		}

		return c.Status(resp.StatusCode()).Send(resp.Body())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
