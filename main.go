package main

import (
	. "../command-line-quiz/models"
	"encoding/json"
	"fmt"
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
			ratingList := getRatingList()
			fmt.Println(getTopFiveRatings(ratingList))
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

func getRandomExpression(expressions []Expression) Expression {
	lastIndex := len(expressions) - 1
	return expressions[rand.Intn(lastIndex)]
}

func askQuestion(isCorrect chan<- bool, done chan bool) {
	expressions := getExpressions()
	go func() {
		for {
			expression := getRandomExpression(expressions)
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
func getExpressions() (expressions []Expression) {
	jsonFile, err := os.Open(expressionFile)
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(jsonByteValue, &expressions)
	if err != nil {
		panic(err)
	}
	return
}

func getRatingList() (ratings []Rating, err error) {
	jsonFile, err := os.Open(ratingFile)
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(jsonByteValue, &ratings)
	return
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
func addToRating(name string, correct int) error {
	previousRating := getRatingList()
	previousRating = append(previousRating, Rating{Name: name, Correct: correct})
	resultRating, err := json.Marshal(previousRating)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(ratingFile)
	_, err = file.WriteString(string(resultRating))
	if err != nil {
		panic(err)
	}
	file.Close()
	fmt.Printf("%s has been added to TOP!\n", name)
	return nil
}
func isTopFive(userCorrectAnswers int) bool {
	ratingList := getRatingList()
	if len(ratingList) < 5 {
		return true
	}
	topFive := getTopFiveRatings(ratingList)
	for _, rating := range topFive {
		if userCorrectAnswers > rating.Correct {
			return true
		}
	}
	return false
}
