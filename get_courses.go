package main

import (
	"fmt"
	"log"
	"time"
)

func getCourses(setup *Setup, prompt string) ([]Course, error) {
	start := time.Now()
	filterPrompt := `
		Extract explicit information from the prompt to create a course filter. Do not infer any details; only include information that is clearly stated in the prompt.
		here are some examples:
		"Can I learn violin this semester?", the course filter would be "Title: violin".
		"I want to take a class with Professor Smith.", the course filter would be "InstructorName: Smith".
		"Show me all the courses that are on Monday.", the course filter would be "Days: Monday".
		"Are there any courses that are online?", the course filter would be "Location: Online".
		"Where does electrical engineering take place?", the course filter would be "Title: Electrical Engineering".
		`
	courseFilter := setup.openAIClient.GetCourseFilter(prompt, filterPrompt)
	log.Printf("Original %s", courseFilter)

	correctedFilter := setup.chromaDB.correctCourseFilter(setup.collections, courseFilter)
	log.Printf("Corrected %s", correctedFilter)

	filteredCourses, err := filterCourses(setup.sqlDB, correctedFilter)
	if err != nil {
		fmt.Println("Error filtering courses:", err)
		return nil, err
	}

	log.Printf("Found %d courses in %.2f seconds", len(filteredCourses), time.Since(start).Seconds())

	return filteredCourses, nil
}
