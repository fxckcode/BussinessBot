package functions

import (
	"github.com/google/generative-ai-go/genai"
)

var GetAllTasks = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "GetAllTasks",
		Description: "Trae todas las tareas",
	}},
}
