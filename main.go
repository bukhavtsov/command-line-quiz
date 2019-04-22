package main

import (
	"encoding/json"
	"fmt"
	. "github.com/bukhavtsov/command-line-quiz/models"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"time"
)

const (
	expressionFile = "jsons/expressions.json"
	ratingFile     = "jsons/rating.json"
)

var expressions []Expression
var ratings []Rating

func init() {
	jsonFile, err := os.Open(expressionFile)
	if err != nil {
		fmt.Printf("expression file `%s` not found", expressionFile)
		os.Exit(1)
	}
	defer jsonFile.Close()
	jsonByteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("failed reading file: %s", err)
		os.Exit(1)
	}
	err = json.Unmarshal(jsonByteValue, &expressions)
	if err != nil {
		fmt.Printf("failed reading JSON: %s", err)
		os.Exit(1)
	}

	jsonFile, err = os.Open(ratingFile)
	if err != nil {
		jsonFile, err = os.Create(expressionFile)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Printf("expression file `%s` has been created", expressionFile)
	}
	jsonByteValue, _ = ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(jsonByteValue, &ratings)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

func main() {
	for {
		fmt.Printf("------------------\n" +
			"To make a choice:\n" +
			"1. Start game\n" +
			"2. Get rating\n" +
			"3. Exit\n" +
			"------------------\n")
		var funcNumber int
		fmt.Scanf("%d\n", &funcNumber)
		switch funcNumber {
		case 1:
			fmt.Println("Let's go!")
			startGame()
		case 2:
			fmt.Println(getTopFiveRatings(ratings))
		case 3:
			os.Exit(1)
		default:
			fmt.Println("Incorrect choice")
		}
	}
}

func startGame() {
	var correct int
	isCorrect := make(chan bool)
	done := make(chan bool)
	doneFunc := make(chan bool)
	until := time.After(time.Minute)
	go func() {
		go askQuestion(isCorrect, done)
		for {
			select {
			case isCorrect := <-isCorrect:
				if isCorrect {
					correct++
				}
			case <-until:
				doneFunc <- true
				done <- true
				return
			}
		}
	}()
	<-doneFunc
	if isTopFive(correct) {
		fmt.Println("\nEnter your name:")
		var name string
		fmt.Scanf("%s\n", &name)
		addToRating(name, correct)
	}
	fmt.Printf("\n----------------\n"+
		"Final result:\n"+
		"number correct answers = %d\n"+
		"----------------\n", correct)
}

func getRandomExpression() Expression {
	lastIndex := len(expressions) - 1
	return expressions[rand.Intn(lastIndex)]
}

func askQuestion(isCorrect chan<- bool, done chan bool) {
	go func() {
		for {
			expression := getRandomExpression()
			var userAnswer string
			fmt.Print(expression.Value, "=")
			fmt.Scanf("%s\n", &userAnswer)
			if expression.Answer == userAnswer {
				isCorrect <- true
			} else {
				isCorrect <- false
			}
		}
	}()
	for {
		select {
		case <-done:
			return
		}
	}
}
func getTopFiveRatings(ratings []Rating) (topFive []Rating) {
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].Correct > ratings[j].Correct
	})
	if len(ratings) >= 5 {
		topFive = ratings[0:5]
		return
	} else {
		return ratings
	}
}
func addToRating(name string, correct int) {
	ratings = append(ratings, Rating{Name: name, Correct: correct})
	resultRating, err := json.Marshal(ratings)
	if err != nil {
		fmt.Println("error:", err)
	}
	file, err := os.Create(ratingFile)
	defer file.Close()
	_, err = file.WriteString(string(resultRating))
	if err != nil {
		fmt.Printf("%s has't been added to TOP!\n", name)
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%s has been added to TOP!\n", name)
}
func isTopFive(userCorrectAnswers int) bool {
	if len(ratings) < 5 {
		return true
	}
	topFive := getTopFiveRatings(ratings)
	for _, rating := range topFive {
		if userCorrectAnswers > rating.Correct {
			return true
		}
	}
	return false
}
