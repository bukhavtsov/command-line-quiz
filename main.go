package main

import (
	. "../command-line-quiz/functions"
	"fmt"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for {
		fmt.Printf("To make a choice:\n" +
			"1. Start game\n" +
			"2. Get rating\n" +
			"3. Exit\n")
		var funcNumber int
		_, err := fmt.Fscan(os.Stdin, &funcNumber)
		if err != nil {
			panic(err)
		}
		switch funcNumber {
		case 1:
			wg.Add(1)
			StartGame()
			wg.Done()
		case 2:
			ratingList := GetRatingList()
			fmt.Println(GetTopFiveRatings(ratingList))
		case 3:
			os.Exit(1)
		default:
			fmt.Println("Incorrect choice")
		}
	}
}
