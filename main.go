package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type TODO struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Go + React Application")
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"Message": "Server is running!"})
	})

	log.Fatal(app.Listen(":8001"))
}
