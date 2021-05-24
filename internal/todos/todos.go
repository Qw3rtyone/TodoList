package todos

import (
	"log"

	database "github.com/Aswin/TodoList/internal/pkg/db/migrations/mysql"
)

type Todo struct {
	ID    int
	Title string
	Note  string
	State bool
}

func (todo Todo) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Todolist(Title,Note) VALUES(?,?)")

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(todo.Title, todo.Note)

	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Inserted!!")
	return id
}
