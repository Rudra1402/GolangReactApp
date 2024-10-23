package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type TODO struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Go + React Application")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	PORT := os.Getenv("PORT")

	todos := []TODO{}

	// Homepage
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"Message": "Server is running!"})
	})

	// Get all TODOs
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Create a TODO
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &TODO{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is required!"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// Update a TODO
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return err
		}

		for _, todo := range todos {
			if todo.ID == id {
				todo.Completed = true
				return c.Status(201).JSON(todo)
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "TODO not found!"})
	})

	// Delete a TODO
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return err
		}

		for index, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:index], todos[index+1:]...)
				return c.Status(201).JSON(fiber.Map{"Message": "Todo deleted!"})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "TODO not found!"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
