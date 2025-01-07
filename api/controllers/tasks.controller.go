package controllers

import (
	"github.com/fxckcode/BussinessBot/api/models"
	"github.com/gofiber/fiber/v2"
)

func GetTasksHandler(c *fiber.Ctx) error {
	task := models.Task{}
	tasks, err := task.FindAllTasks()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(tasks)
}

func GetTaskHandler(c *fiber.Ctx) error {
	task := models.Task{}
	id := c.Params("id")
	t, err := task.FindTaskById(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(t)
}