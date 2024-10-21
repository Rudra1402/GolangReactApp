package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Go + React Application")
	app := fiber.New()
	log.Fatal(app.Listen(":8001"))
}
