package graph

import "github.com/Aswin/TodoList/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

func NewResolver() *Resolver {

	users := make([]*model.User, 0)

	users = append(users, &model.User{ID: "1", Name: "fphilip"})

	users = append(users, &model.User{ID: "2", Name: "lturanga"})

	return &Resolver{

		todos: make([]*model.Todo, 0),

		users: users,
	}

}

type Resolver struct {
	todos []*model.Todo

	users []*model.User

	lastTodoId int
}
