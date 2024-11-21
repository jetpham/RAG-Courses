package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go"
)

func toolCallingAgent(setup Setup, prompt string) string {
	// using the official example:
	//https://github.com/openai/openai-go/blob/main/examples/chat-completion-tool-calling/main.go

	systemPrompt := `
		You are a Retrieval Augmented Generation model that assists university students using course information.
	`
	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		}),
		Tools: openai.F([]openai.ChatCompletionToolParam{
			{
				Type: openai.F(openai.ChatCompletionToolTypeFunction),
				Function: openai.F(openai.FunctionDefinitionParam{
					Name:        openai.String("get_courses"),
					Description: openai.String("get course information"),
					Parameters: openai.F(openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"prompt": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"prompt"},
					}),
				}),
			},
		}),
		Model: openai.F(openai.ChatModelGPT4oMini),
	}

	completion, err := setup.openAIClient.client.Chat.Completions.New(context.TODO(), params)
	if err != nil {
		log.Printf("Error creating chat completion: %v", err)
		return ""
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// If there was not tool calls, crashout
	if len(toolCalls) == 0 {
		log.Printf("No function call")
		return completion.Choices[0].Message.Content
	} else {
		log.Printf("Function call: %v\n", toolCalls[0].Function.Name)
		fmt.Printf("Function call: %v\n", toolCalls[0].Function.Name)
	}

	// If there was tool calls, continue
	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_courses" {
			// Extract the prompt from the function call arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				log.Printf("Error unmarshalling arguments: %v", err)
				continue
			}
			prompt := args["prompt"].(string)

			log.Printf("%v(\"%s\")", toolCalls[0].Function.Name, prompt)

			// Call the getCourses function with the arguments requested by the model
			courses, err := getCourses(&setup, prompt)
			if err != nil {
				log.Printf("Error getting courses: %v", err)
				continue
			}
			coursesJSON, _ := json.Marshal(courses)
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, string(coursesJSON)))
		}
	}

	params.Messages.Value = append(params.Messages.Value, openai.SystemMessage(`
		Task: Generate an answer that corresponds to the provided question, mimicking the question's structure and format. Ensure the response is succinct, directly relevant to the query, and excludes any extraneous details.

		Reponce Format:
			[answer to the question]
			cited courses:
			[course details]
			
		Course Details: If relevant courses are involved in the answer, format them as follows:
		- Format: (subject_code)(course_number)-(section) title_short_desc by primary_instructor_full_name (relevant course details).
		- Ensure each course listed accurately responds to the original question.

		Response Guidelines:
		- Provide responses in plain text format, avoiding markdown.
		- List course details in a numbered format for clarity.
		- Ensure responce directory addresses the question
		`))
	params.Messages.Value = append(params.Messages.Value, openai.UserMessage(fmt.Sprintf("Prompt: %s", prompt)))
	completion, err = setup.openAIClient.client.Chat.Completions.New(context.TODO(), params)
	if err != nil {
		log.Printf("Error creating chat completion: %v", err)
		return ""
	}

	return completion.Choices[0].Message.Content
}
