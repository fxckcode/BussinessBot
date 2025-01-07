package tools

import "github.com/fxckcode/BussinessBot/api/models"

func GetAllTasks() (*[]models.Task, string) {
	task := models.Task{}
	tasks, err := task.FindAllTasks()
	if err != nil {
		return nil, "Se ha producido un error"
	}

	return tasks, ""
}