package todos

import (
	"database/sql"
	"log"
	"strconv"

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

func (todo Todo) UpdateTodo(id string) (int64, error) {
	stmt, err := database.Db.Prepare("UPDATE `Todolist` SET `Title` = ?, `Note` = ?, `Completed` = ? WHERE `ID` = ?")

	if err != nil {
		log.Fatal(err)
	}

	updateID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(todo.Title, todo.Note, todo.State, updateID)
	if err != nil {
		log.Fatal(err)
	}

	response, err := res.RowsAffected()
	if err != nil {
		log.Print(err)
	}
	log.Print("Update finished")

	return response, err
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

func GetById(searchID string) (Todo, error) {

	stmt, err := database.Db.Prepare("Select * FROM `Todolist` WHERE `ID` = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	searchIDInt, err := strconv.Atoi(searchID)
	if err != nil {
		log.Fatal(err)
	}

	row := stmt.QueryRow(searchIDInt)

	var todo Todo

	err = row.Scan(&todo.ID, &todo.Title, &todo.Note, &todo.State)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}

		return Todo{}, err
	}

	return todo, nil
}

func GetPortion(state bool) []Todo {
	stmt, err := database.Db.Prepare("SELECT * FROM `Todolist` WHERE `Completed` = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(state)
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
