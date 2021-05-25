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
	stmt, err := database.Db.Prepare("INSERT INTO Todolist(Title,Note,Completed) VALUES(?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(todo.Title, todo.Note, todo.State)

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

func GetAll() []Todo {
	stmt, err := database.Db.Prepare("SELECT * FROM `Todolist`")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Note, &todo.State)

		if err != nil {
			log.Fatal(err)
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return todos
}
