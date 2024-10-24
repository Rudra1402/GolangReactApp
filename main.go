package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TODO struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Go + React Application")
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	PORT := os.Getenv("PORT")
	MONGODB_URI := os.Getenv("MONGODB_URI")

	if PORT == "" {
		PORT = "8081"
	}

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB created successfully!")

	collection = client.Database("golang_db").Collection("todos")

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", addTodos)
	// app.Patch("/api/todos/:id", updateTodos)
	// app.Delete("/api/todos/:id", deleteTodos)

	// todos := []TODO{}

	// Homepage
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"Message": "Server is running!"})
	})

	// // Get all TODOs
	// app.Get("/api/todos", func(c *fiber.Ctx) error {
	// 	return c.Status(200).JSON(todos)
	// })

	// // Create a TODO
	// app.Post("/api/todos", func(c *fiber.Ctx) error {
	// 	todo := &TODO{}

	// 	if err := c.BodyParser(todo); err != nil {
	// 		return err
	// 	}

	// 	if todo.Body == "" {
	// 		return c.Status(400).JSON(fiber.Map{"error": "Body is required!"})
	// 	}

	// 	todo.ID = len(todos) + 1
	// 	todos = append(todos, *todo)

	// 	return c.Status(201).JSON(todo)
	// })

	// // Update a TODO
	// app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
	// 	id, err := c.ParamsInt("id")

	// 	if err != nil {
	// 		return err
	// 	}

	// 	for _, todo := range todos {
	// 		if todo.ID == id {
	// 			todo.Completed = true
	// 			return c.Status(201).JSON(todo)
	// 		}
	// 	}

	// 	return c.Status(404).JSON(fiber.Map{"error": "TODO not found!"})
	// })

	// // Delete a TODO
	// app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
	// 	id, err := c.ParamsInt("id")

	// 	if err != nil {
	// 		return err
	// 	}

	// 	for index, todo := range todos {
	// 		if todo.ID == id {
	// 			todos = append(todos[:index], todos[index+1:]...)
	// 			return c.Status(201).JSON(fiber.Map{"Message": "Todo deleted!"})
	// 		}
	// 	}

	// 	return c.Status(404).JSON(fiber.Map{"error": "TODO not found!"})
	// })

	log.Fatal(app.Listen(":" + PORT))
}

func getTodos(c *fiber.Ctx) error {
	var todos []TODO

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo TODO
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.Status(200).JSON(todos)
}

func addTodos(c *fiber.Ctx) error {
	todo := new(TODO)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Body is required!"})
	}

	insertTodo, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertTodo.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

// func updateTodos(c *fiber.Ctx) error {}
// func deleteTodos(c *fiber.Ctx) error {}
