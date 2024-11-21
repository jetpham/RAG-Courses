package main

import (
	"context"
	"fmt"
	"log"
	"os"

	chroma "github.com/amikos-tech/chroma-go"
	openai "github.com/amikos-tech/chroma-go/openai"
)

type chromaDB struct {
	ctx      context.Context
	openaiEf *openai.OpenAIEmbeddingFunction
	client   *chroma.Client
}

func newChroma() (*chromaDB, error) {
	defer log.Println("Chroma setup complete")
	ctx := context.Background()

	// Set up chroma-go client
	client, err := chroma.NewClient("http://localhost:8000")
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	// Reset the client to clear any previous state
	// client.Reset(ctx)

	openaiEf, err := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenAI embedding function: %s", err)
	}

	return &chromaDB{
		ctx:      ctx,
		openaiEf: openaiEf,
		client:   client,
	}, nil
}
