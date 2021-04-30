package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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
	ioutil.WriteFile("storage/todoList.json", file, 0644)
	fmt.Println("Finished writing")
}
func showList(list TodoList) {
	fmt.Println("Trying to print list. Length " + strconv.Itoa(len(list.TodoList)))
	for i := 0; i < len(list.TodoList); i++ {
		fmt.Println("---------------------------------------")
		fmt.Println("ID: " + strconv.Itoa(list.TodoList[i].Id))
		fmt.Println("Title: " + list.TodoList[i].Title)
		fmt.Println("Note: " + list.TodoList[i].Note)
		fmt.Println("Due: " + list.TodoList[i].Due.String())
		fmt.Println("---------------------------------------")

	}
}
func readList() TodoList {
	data, err := os.Open("storage/todoList.json")
	checkerr(err)
	fmt.Println("Opened List!")

	defer data.Close()

	byteVal, _ := ioutil.ReadAll(data)

	var tdList TodoList

	json.Unmarshal(byteVal, &tdList)

	showList(tdList)

	return tdList
}
func printItem(item Todo, w http.ResponseWriter) {
	fmt.Fprintln(w, "---------------------------------------")
	fmt.Fprintln(w, "ID: "+strconv.Itoa(item.Id))
	fmt.Fprintln(w, "Title: "+item.Title)
	fmt.Fprintln(w, "Note: "+item.Note)
	fmt.Fprintln(w, "Due: "+item.Due.String())
	fmt.Fprintln(w, "---------------------------------------")
}
func printList(list TodoList, w http.ResponseWriter) {
	fmt.Println("Trying to print list. Length " + strconv.Itoa(len(list.TodoList)))
	for i := 0; i < len(list.TodoList); i++ {
		printItem(list.TodoList[i], w)

	}
}
func homePage(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 50; i++ {
		fmt.Fprintln(w, "Welcome, please work <3")
	}
	fmt.Println("End of page")
}
func getFullList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: full list")
	printList(readList(), w)
}
func getOneItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: single item")
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Fprintln(w, "Key sent: "+key)

	for _, item := range readList().TodoList {
		k, err := strconv.Atoi(key)
		fmt.Println(k)
		if (err == nil) && (k == item.Id) {
			printItem(item, w)
		}
	}
}

/*func addToList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: Add new item")

}
*/
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/showAll", getFullList)
	myRouter.HandleFunc("/show/{id}", getOneItem)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	if _, err := os.Stat("storage/todoList.json"); os.IsNotExist(err) {
		initStore()
	}

	readList()

	handleRequests()

}
