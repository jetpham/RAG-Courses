package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Setup struct {
	courses      []Course
	openAIClient *OpenAIClient
	chromaDB     *chromaDB
	sqlDB        *gorm.DB
	collections  *collections
}

func newSetup() (Setup, error) {
	startTime := time.Now()
	defer func() {
		log.Printf("Setup complete in %.2fs", time.Since(startTime).Seconds())
	}()

	if os.Getenv("OPENAI_API_KEY") == "" {
		log.Println("Loading OPENAI_API_KEY from .env file")
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading .env file:", err)
			return Setup{}, err
		}
	} else {
		log.Println("Using OPENAI_API_KEY from environment")
	}

	// prevent timeouts with openai
	os.Setenv("MODULES_CLIENT_TIMEOUT", "2m")

	courses, err := loadCSV("Fall 2024 Class Schedule.csv")
	if err != nil {
		fmt.Println("Error loading CSV file:", err)
		return Setup{}, err
	}

	openAIClient := NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))
	if openAIClient == nil {
		fmt.Println("Error creating OpenAI client")
		return Setup{}, err
	}

	chromaDB, err := newChroma()
	if err != nil {
		fmt.Println("Error setting up ChromaDB:", err)
		return Setup{}, err
	}

	sqlDB, err := newSqlite(courses)
	if err != nil {
		fmt.Println("Error setting up SQLite database:", err)
		return Setup{}, err
	}

	collections, err := makeCollections(chromaDB, courses)
	if err != nil {
		fmt.Println("Error loading collections:", err)
		return Setup{}, err
	}

	return Setup{
		courses:      courses,
		openAIClient: openAIClient,
		chromaDB:     chromaDB,
		sqlDB:        sqlDB,
		collections:  collections,
	}, nil
}
