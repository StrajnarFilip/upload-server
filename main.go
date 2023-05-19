package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Post("/upload", func(c *fiber.Ctx) error {
		// Parse the multipart form:
		if form, err := c.MultipartForm(); err == nil {
			// Get all files from "documents" key:
			files := form.File["documents"]

			// Loop through files:
			for _, file := range files {
				// Create a random buffer with 32 bytes
				buffer := make([]byte, 32)
				rand.Read(buffer)

				// Hexadecimal representation of random buffer
				hexString := hex.EncodeToString(buffer)

				// File path
				filePath := fmt.Sprintf("/%s%s", hexString, file.Filename)

				// Save the files to disk:
				if err := c.SaveFile(file, "./public"+filePath); err != nil {
					return err
				}
				return c.SendString(filePath)
			}
			return err
		}
		return nil
	})
	app.Static("/", "./public")
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "*",
		AllowMethods:     "*",
		AllowCredentials: true,
	}))

	address, defined := os.LookupEnv("UPLOADSERVERADDRESS")
	if defined {
		app.Listen(address)
	} else {
		app.Listen("127.0.0.1:8080")
	}
}
