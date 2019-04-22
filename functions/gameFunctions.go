package functions

import (
	. "../models"
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

func StartGame() {
	var correct int
	var incorrect int
	go calculateExpressions(&correct, &incorrect)
	time.Sleep(time.Minute)
	if isTopFive(correct) {
		fmt.Println("Enter your name:")
		var name string
		fmt.Fscan(os.Stdin, &name)
		addToRating(name, correct)
	}
	fmt.Printf("Final result:\n"+
		"number correct answers = %d\n"+
		"number Incorrect answers = %d", correct, incorrect)
}

func calculateExpressions(numberCorrectAnswers *int, numberIncorrectAnswers *int) {
	expressions := getExpressions()
	shuffle(expressions)
	for _, expression := range expressions {
		var userAnswer string
		fmt.Println(expression.Value)
		fmt.Fscan(os.Stdin, &userAnswer)
		if expression.Answer == userAnswer {
			*numberCorrectAnswers++
		} else {
			*numberIncorrectAnswers++
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

func GetRatingList() (ratings []Rating) {
	jsonFile, err := os.Open(ratingFile)
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	jsonByteValue, _ := ioutil.ReadAll(jsonFile)
	_ = json.Unmarshal(jsonByteValue, &ratings)
	return
}
func GetTopFiveRatings(ratings []Rating) (topFive []Rating) {
	sort.Slice(ratings, func(i, j int) bool {
		return ratings[i].NumberCorrectAnswers > ratings[j].NumberCorrectAnswers
	})
	if len(ratings) >= 5 {
		topFive = ratings[0:5]
		return
	} else {
		return ratings
	}
}
func addToRating(name string, numberCorrectAnswers int) {
	previousRating := GetRatingList()
	previousRating = append(previousRating, Rating{Name: name, NumberCorrectAnswers: numberCorrectAnswers})
	resultRating, err := json.Marshal(previousRating)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(ratingFile)
	file.WriteString(string(resultRating))
	defer file.Close()
	fmt.Println("You has been added to TOP!")
}
func isTopFive(numberUserCorrectAnswers int) bool {
	ratingList := GetRatingList()
	topFive := GetTopFiveRatings(ratingList)
	for _, rating := range topFive {
		if numberUserCorrectAnswers >= rating.NumberCorrectAnswers {
			return true
		}
	}
	return false
}
func shuffle(expressions []Expression) {
	rand.Shuffle(len(expressions), func(i, j int) {
		expressions[i], expressions[j] = expressions[j], expressions[i]
	})
}
