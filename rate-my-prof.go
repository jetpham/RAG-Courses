package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ProfessorInfo struct {
	Name           string      `json:"name"`
	Department     string      `json:"department"`
	School         string      `json:"school"`
	Rating         float64     `json:"rating"`
	Difficulty     float64     `json:"difficulty"`
	TotalRatings   int         `json:"total_ratings"`
	WouldTakeAgain interface{} `json:"would_take_again"`
}

// make a get request to the flask server running the api to scrape rate my professor
func getRateMyProfessorData(name string) (ProfessorInfo, error) {
	url := "http://localhost:5000/professor?name=" + strings.ReplaceAll(name, " ", "+")
	resp, err := http.Get(url)
	if err != nil {
		return ProfessorInfo{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ProfessorInfo{}, err
	}

	var professorInfo ProfessorInfo
	err = json.Unmarshal(body, &professorInfo)
	if err != nil {
		fmt.Println(err)
		return ProfessorInfo{}, err
	}

	return professorInfo, nil
}
