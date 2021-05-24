package graph

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/Aswin/TodoList/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type TodoList struct {
	TodoList []Todo `json:"todolist"`
}

type Todo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Note  string `json:"note"`
	State bool   `json:"state"`
}

func readList() TodoList {
	data, err := os.Open("storage/todoList.json")
	if err != nil {
		log.Fatal("Error: ", err)
	}

	//checkerr(err, false)
	//fmt.Println("Opened file!!")

	defer data.Close()

	byteVal, _ := ioutil.ReadAll(data)

	var tdList TodoList

	json.Unmarshal(byteVal, &tdList)

	return tdList
}
func NewResolver() *Resolver {

	listtd := readList()

	todos := make([]*model.Todo, 0)

	for _, item := range listtd.TodoList {
		todos = append(todos, &model.Todo{ID: item.Id, Title: item.Title, Text: item.Note, Done: item.State})
	}

	return &Resolver{

		todos: todos,

		lastTodoId: len(listtd.TodoList),
	}

}

type Resolver struct {
	todos []*model.Todo

	lastTodoId int
}
