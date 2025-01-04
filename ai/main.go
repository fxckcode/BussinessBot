package ai

import (
	"context"
	"fmt"
	"log"

	"github.com/fxckcode/BussinessBot/ai/config"
	"github.com/fxckcode/BussinessBot/env"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// createClient crea un nuevo cliente de genai.
func createClient(token string) *genai.Client {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Verificamos si el cliente es nil
	if client == nil {
		log.Fatal("Error: No se pudo crear el cliente")
	}

	return client
}

func SearchGemini(ctx context.Context,content string, modelAI string) string {
	token := env.ViperEnvVariable("GEMINI_API_KEY")

	client := createClient(token)
	model := client.GenerativeModel(modelAI)

	// Config
	model.SetTemperature(float32(config.TEMPERATURE))
	model.SystemInstruction = genai.NewUserContent(genai.Text(config.PROMPT))

	response, err := model.GenerateContent(ctx, genai.Text(content))

	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	output := make(chan string)
	go outputResponse(response, output)

	return <-output
}

func outputResponse(resp *genai.GenerateContentResponse, output chan string) {
	if resp != nil && len(resp.Candidates) > 0 {
		firstCandidate := resp.Candidates[0]
		if firstCandidate.Content != nil && len(firstCandidate.Content.Parts) > 0 {
			part := fmt.Sprint(firstCandidate.Content.Parts[0])
			output <- part
		} else {
			output <- "no content in response"
		}
	} else {
		output <- "response is empty"
	}
	close(output)
}

