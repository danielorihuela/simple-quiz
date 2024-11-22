package models

type Question struct {
	ID            int
	Title         string
	Options       []string
	CorrectAnswer int
}
