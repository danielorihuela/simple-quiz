package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danielorihuela/simple-quizz/models"
)

var questions = []models.Question{
	{
		ID:            1,
		Title:         "What's the result of the operation \"4 + 4\"?",
		Options:       []string{"7", "8", "9"},
		CorrectAnswer: 1,
	},
	{
		ID:            2,
		Title:         "What's the result of the operation \"4 * 4\"?",
		Options:       []string{"16", "17", "18", "19"},
		CorrectAnswer: 0,
	},
	{
		ID:            3,
		Title:         "What's the result of the operation \"4 / 4\"?",
		Options:       []string{"3", "2", "1"},
		CorrectAnswer: 2,
	},
}

var results []int

func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func postAnswers(w http.ResponseWriter, r *http.Request) {
	var answers []int
	err := json.NewDecoder(r.Body).Decode(&answers)
	if err != nil {
		log.Println("POST ANSWERS: " + err.Error())
		http.Error(w, "Invalid data", 400)
	}

	correctAnswers := 0
	for i, q := range questions {
		if answers[i] == q.CorrectAnswer {
			correctAnswers++
		}
	}

	resultAbove := 0
	for _, result := range results {
		if correctAnswers > result {
			resultAbove++
		}
	}

	var betterThan float32 = 100
	if len(results) > 0 {
		betterThan = (float32(resultAbove) / float32(len(results))) * 100
	}

	response := models.PostAnswers{
		Result:     correctAnswers,
		BetterThan: betterThan,
	}

	results = append(results, correctAnswers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/questions", getQuestions)
	http.HandleFunc("/answers", postAnswers)
	http.ListenAndServe(":8080", nil)
}
