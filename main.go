package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"
)

func main() {
	// Load templates (Not mandatory)
	engine := html.New("./views", ".html")

	app := *fiber.New(fiber.Config{
		Views: engine,
	})

	// Use built in logger
	app.Use(logger.New())

	// Fiber does not recover from panic by default.
	// Use the following middleware to recover from panic.
	app.Use(recover.New())

	// Basic Route and middleware
	app.Get("/",
		MyCustomMiddleWare,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "Welcome to fiber in golang"})
		},
	)

	// Group Route
	myGroup := app.Group("/api", MyCustomMiddleWare)
	myGroup.Get("/one", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello form api one",
		})
	})
	myGroup.Get("/two", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello form api two",
		})
	})

	// Static file access
	app.Static("/", "./public")

	// Render a template
	app.Get("/template", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"name": "John Doe",
		})
	})

	app.Post("/template", func(c *fiber.Ctx) error {
		var body struct {
			Message string
		}

		err := c.BodyParser(&body)
		if err != nil {
			return err
		}

		return c.Render("index", fiber.Map{
			"name":    "John Snow",
			"message": body.Message,
		})
	})

	// Routing with parameters
	app.Get("/names/:name/age/:age", func(c *fiber.Ctx) error {
		name := c.Params("name")
		age, _ := strconv.Atoi(c.Params("age"))

		c.SendStatus(http.StatusAccepted)
		return c.JSON(fiber.Map{
			"status": "success",
			"name":   name,
			"age":    age,
		})
	})

	log.Fatal(app.Listen(":4000"))
}

// Custom middleware/handler
func MyCustomMiddleWare(c *fiber.Ctx) error {
	if 1 == 0 {
		c.SendStatus(http.StatusForbidden)
		return c.JSON(fiber.Map{"message": "You have been denied!"})
	}
	return c.Next()
}
