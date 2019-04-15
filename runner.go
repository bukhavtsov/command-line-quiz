package main

import (
	. "../command-line-quiz/functions"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("To make a choice:\n" +
		"1. Start game\n" +
		"2. Get rating\n")
	var funcNumber int
	_, err := fmt.Fscan(os.Stdin, &funcNumber)
	if err != nil {
		panic(err)
	}
	switch funcNumber {
	case 1:
		StartGame()
	case 2:
		ratingList := GetRatingList()
		fmt.Println(GetTopFiveRatings(ratingList))
	default:
		fmt.Println("Incorrect choice")
	}
}
