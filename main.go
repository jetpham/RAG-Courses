package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Println("setting up...")
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	setup, err := newSetup()
	if err != nil {
		fmt.Println("Error during setup:", err)
		return
	}
	fmt.Print("\033[H\033[2J")
	fmt.Println("Ask a question about courses or type 'exit' to quit.")
	exampleQuestions := []string{
		"What courses is Phil Peterson teaching in Fall 2024?",
		"Which philosophy courses are offered this semester?",
		"Where does Bioinformatics meet?",
		"Can I learn guitar this semester?",
		"I would like to take a Rhetoric course from Phil Choong. What can I take?",
	}

	fmt.Println("Example questions:")
	for _, question := range exampleQuestions {
		fmt.Println("-", question)
	}
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter prompt: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)
		if prompt == "exit" {
			break
		}
		output := toolCallingAgent(setup, prompt)
		fmt.Println(output)
	}
}
