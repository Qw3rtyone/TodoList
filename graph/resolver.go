package graph

import (
	"github.com/Aswin/TodoList/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver() *Resolver {
	/*
		todos := make([]*model.Todo, 0)

		todos = append(todos, &model.Todo{ID: "1", Title: "First Todo", Text: "Body of todo", Done: false})
		todos = append(todos, &model.Todo{ID: "2", Title: "Get milk", Text: "Milk a cow", Done: true})

		return &Resolver{

			todos: todos,

			lastTodoId: 3,
		}
	*/

	return &Resolver{}
}

type Resolver struct {
	todos []*model.Todo

	lastTodoId int
}
