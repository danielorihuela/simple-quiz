package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/danielorihuela/simple-quizz/models"
	"github.com/spf13/cobra"
)

var scheme = "http"
var host = "localhost:8080"
var questionsUrl = url.URL{
	Scheme: scheme,
	Host:   host,
	Path:   "/questions",
}

var answersUrl = url.URL{
	Scheme: scheme,
	Host:   host,
	Path:   "/answers",
}

var rootCmd = &cobra.Command{
	Use:   "simple-quiz-cli",
	Short: "Test your maths knowledge with a new quiz every day!",
	Long: `Get a new math quiz every day and compete with other users
around the world.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(questionsUrl.String())
		exitIfError(err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		exitIfError(err)

		var questions []models.Question
		err = json.Unmarshal(body, &questions)
		exitIfError(err)

		var answers []int
		for _, q := range questions {
			fmt.Printf("Q%d: %s\n", q.ID, q.Title)
			for i, opt := range q.Options {
				fmt.Printf("%c. %s\n", 'a'+i, opt)
			}

			var answer = readUserAnswer(len(q.Options))
			answers = append(answers, answer)
		}

		jsonValue, err := json.Marshal(answers)
		exitIfError(err)

		fmt.Println("Submitting your answers...")
		resp, err = http.Post(answersUrl.String(), "application/json", bytes.NewBuffer(jsonValue))
		exitIfError(err)
		defer resp.Body.Close()

		var result models.PostAnswers
		json.NewDecoder(resp.Body).Decode(&result)
		if result.Result == 1 {
			fmt.Printf("You answered %d question correctly.", result.Result)
		} else {
			fmt.Printf("You answered %d questions correctly.", result.Result)
		}
		fmt.Printf("\nYou were better than %.0f%% of all quizzers.", result.BetterThan)
	},
}

func exitIfError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func readUserAnswer(numOptions int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Which is your answer? ")
		input, err := reader.ReadString('\n')
		exitIfError(err)

		lastOptionValue := 'a' + numOptions - 1
		if len(input) == 2 && input[1] == '\n' && input[0] >= 'a' && input[0] <= byte(lastOptionValue) {
			fmt.Println()
			return int(input[0]) - 'a'
		}

		fmt.Printf("Invalid input. Please enter a character between \"a\" and \"%c\".\n", lastOptionValue)
	}
}

func main() {
	err := rootCmd.Execute()
	exitIfError(err)
}
