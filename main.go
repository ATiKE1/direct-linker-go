package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New()

	UPLOAD_DIR := "uploads"
	os.MkdirAll(UPLOAD_DIR, os.ModePerm)

	app.Static("/uploads", UPLOAD_DIR)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "File not attached",
			})
		}

		fileID := uuid.New().String()
		ext := filepath.Ext(fileHeader.Filename)
		filename := fileID + ext
		filePath := filepath.Join(UPLOAD_DIR, filename)

		if err := c.SaveFile(fileHeader, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot save file",
			})
		}

		downloadURL := fmt.Sprintf("/uploads/%s", filename)
		return c.JSON(fiber.Map{"url": downloadURL})
	})

	fmt.Println("Server is running at http://localhost:5000")
	app.Listen(":5000")
}
