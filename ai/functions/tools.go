package functions

import (
	"github.com/google/generative-ai-go/genai"
)

var PrintHello = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "PrintHello",
		Description: "Print hello world",
	}},
}
