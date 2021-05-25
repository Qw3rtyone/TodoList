package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strconv"

	"github.com/Aswin/TodoList/graph/generated"
	"github.com/Aswin/TodoList/graph/model"
	"github.com/Aswin/TodoList/internal/todos"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	//panic(fmt.Errorf("not implemented"))
	/*
		newTodo := &model.Todo{
			ID:    strconv.Itoa(r.lastTodoId),
			Title: input.Title,
			Text:  input.Text,
			Done:  false,
		}

		r.todos = append(r.todos, newTodo)
		r.lastTodoId++

		return newTodo, nil
	*/

	var todo todos.Todo
	todo.Title = input.Title
	todo.Note = input.Text
	todo.State = false
	todoID := todo.Save()

	return &model.Todo{ID: strconv.FormatInt(todoID, 10), Title: todo.Title, Text: todo.Note, Done: todo.State}, nil

}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	var resTodos []*model.Todo
	var dbTodos []todos.Todo

	dbTodos = todos.GetAll()

	for _, item := range dbTodos {
		resTodos = append(resTodos, &model.Todo{ID: strconv.Itoa(item.ID), Title: item.Title, Text: item.Note, Done: item.State})
	}

	return resTodos, nil
	//return r.todos, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
