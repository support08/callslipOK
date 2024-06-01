package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

var client = &http.Client{}

func uploadFile(c *fiber.Ctx) error {
	branchId := 22449
	apiKey := "SLIPOKWBD62MC"
	filePath := ""

	url := fmt.Sprintf("https://api.slipok.com/api/line/apikey/%d", branchId)

	// Create a buffer to write our multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}
	defer file.Close()

	// Create the file field in the multipart form
	part, err := writer.CreateFormFile("files", filePath)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}
	if _, err = io.Copy(part, file); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}

	// // Add the log field
	// if err = writer.WriteField("log", "true"); err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"code":    "500",
	// 		"message": err.Error(),
	// 	})
	// }

	// if err = writer.Close(); err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"code":    "500",
	// 		"message": err.Error(),
	// 	})
	// }

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}

	// Set the headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-authorization", apiKey)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}

	// Parse the response
	var jsonResponse map[string]interface{}
	if err = json.Unmarshal(responseBody, &jsonResponse); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": err.Error(),
		})
	}

	// Check the success field in the JSON response
	if success, ok := jsonResponse["success"].(bool); ok && success {
		return c.JSON(jsonResponse["data"])
	} else {
		errorCode := jsonResponse["code"].(float64)
		errorMessage := jsonResponse["message"].(string)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    errorCode,
			"message": errorMessage,
		})
	}
}

func main() {
	app := fiber.New()
	app.Post("/upload", uploadFile)
	app.Listen(":8000")
}
