package main

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Hardcode your values directly in the code
	branchID := "22345"
	apiKey := "SLIPOK6H4K8YP"
	fileLocation := "C:/Support08 Work/NOTE CODE/sql tech/export/154210.jpg"

	app := fiber.New()

	app.Post("/call-api", func(c *fiber.Ctx) error {
		client := resty.New()
		resp, err := client.R().
			SetHeader("x-authorization", apiKey).
			SetFile("files", fileLocation).
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

	log.Fatal(app.Listen(":3000"))
}
