package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

type TodoList struct {
	TodoList []Todo `json:"todolist"`
}

type Todo struct {
	Id    int       `json:"id"`
	Title string    `json:"title"`
	Note  string    `json:"note"`
	Due   time.Time `json:"due"`
}

func checkerr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func initStore() {

	data := TodoList{
		TodoList: []Todo{
			{
				Id:    0,
				Title: "Default",
				Note:  "Default Note",
				Due:   time.Now(),
			},
			{
				Id:    1,
				Title: "Add New Todo",
				Note:  "Add a third todo!",
				Due:   time.Now(),
			},
		},
	}

	file, _ := json.MarshalIndent(data, "", "")
	fmt.Println("writing")
	ioutil.WriteFile("test.json", file, 0644)
	fmt.Println("Finished writing")
}
func showList(list TodoList) {
	fmt.Println("Trying to print list. Length " + strconv.Itoa(len(list.TodoList)))
	for i := 0; i < len(list.TodoList); i++ {
		fmt.Println("ID: " + strconv.Itoa(list.TodoList[i].Id))
		fmt.Println("Title: " + list.TodoList[i].Title)
		fmt.Println("Note: " + list.TodoList[i].Note)
		fmt.Println("Due: " + list.TodoList[i].Due.String())
	}
}

func main() {
	if _, err := os.Stat("test.json"); os.IsNotExist(err) {
		initStore()
	}

	data, err := os.Open("test.json")
	checkerr(err)
	fmt.Println("Opened List!")

	defer data.Close()

	byteVal, _ := ioutil.ReadAll(data)

	var tdList TodoList

	json.Unmarshal(byteVal, &tdList)

	showList(tdList)
}
