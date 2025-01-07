package ai

import (
	"context"
	"fmt"
    "log"
    "strings"

	"github.com/fxckcode/BussinessBot/ai/config"
	"github.com/fxckcode/BussinessBot/ai/functions"
	"github.com/fxckcode/BussinessBot/api/controllers/tools"
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

func SearchGemini(ctx context.Context, content string, modelAI string) string {
    token := env.ViperEnvVariable("GEMINI_API_KEY")

    client := createClient(token)
    model := client.GenerativeModel(modelAI)

    // Tools
    model.Tools = []*genai.Tool{functions.GetAllTasks}

    // Config
    model.SetTemperature(float32(config.TEMPERATURE))
    model.SystemInstruction = genai.NewUserContent(genai.Text(config.PROMPT))

    session := model.StartChat()

    resp, err := session.SendMessage(ctx, genai.Text(content))

    if err != nil {
        log.Fatalf("Error generating content: %v", err)
    }

    var response *genai.GenerateContentResponse

    part := resp.Candidates[0].Content.Parts[0]
    funcall, ok := part.(genai.FunctionCall)

    if !ok {
        response, err = model.GenerateContent(ctx, genai.Text(content))
        if err != nil {
            log.Fatalf("Error generating content: %v", err)
        }
    }

    if ok {
        // Handle the function call
        apiResult := handleFunctionCall(funcall)

        // Send the function result back to the generative model
        response, err = session.SendMessage(ctx, genai.FunctionResponse{
            Name:     functions.GetAllTasks.FunctionDeclarations[0].Name,
            Response: convertMap(apiResult),
        })
        if err != nil {
            log.Fatalf("Error sending function response: %v", err)
        }
    }

    output := make(chan string)
    go outputResponse(response, output)

    return <-output
}

func handleFunctionCall(funcall genai.FunctionCall) map[string]string {
    if funcall.Name == "GetAllTasks" {
        tasks, err := tools.GetAllTasks()
        if err != "" {
            log.Printf("Error getting tasks: %v", err)
            return map[string]string{"error": "Failed to get tasks"}
        }

        // Extraer y formatear los datos importantes del JSON
        var taskList []string
        for _, task := range *tasks {
            taskList = append(taskList, fmt.Sprintf("ID: %s, Nombre: %s, Descripción: %s", task.ID, task.Title, task.Description))
        }
        formattedTasks := strings.Join(taskList, "\n")

        return map[string]string{"tasks": formattedTasks}
    }

    return map[string]string{"error": "Unknown function call"}
}

func convertMap(input map[string]string) map[string]any {
	output := make(map[string]any)
	for k, v := range input {
		output[k] = v
	}
	return output
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

