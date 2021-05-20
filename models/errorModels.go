package models

type ErrorList struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
