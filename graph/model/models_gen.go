// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewTodo struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Respose struct {
	Change string `json:"change"`
}

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
}

type UpdateTodo struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Done  bool   `json:"done"`
}
