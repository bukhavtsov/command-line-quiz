package models

type Rating struct {
	Name                 string `json:"name"`
	NumberCorrectAnswers int    `json:"numberCorrectAnswers"`
}
