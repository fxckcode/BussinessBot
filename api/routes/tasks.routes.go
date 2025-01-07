package routes

import (
	"github.com/fxckcode/BussinessBot/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func TasksRoutes(c *fiber.App) {
	route := c.Group("/api/v1")

	route.Get("/tasks", controllers.GetTasksHandler)
	route.Get("/tasks/:id", controllers.GetTaskHandler)
}