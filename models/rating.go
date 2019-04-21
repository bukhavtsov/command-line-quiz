package models

type Rating struct {
	Name    string `json:"name"`
	Correct int    `json:"numberCorrectAnswers"`
}
