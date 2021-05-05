package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type TodoList struct {
	TodoList []Todo `json:"todolist"`
}

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Note  string `json:"note"`
	State bool   `json:"state"`
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
				State: false,
			},
			{
				Id:    1,
				Title: "Add New Todo",
				Note:  "Add a third todo!",
				State: false,
			},
			{
				Id:    2,
				Title: "Add Test Todo",
				Note:  "Add a Fourth todo!",
				State: false,
			},
			{
				Id:    3,
				Title: "Add Next",
				Note:  "Add a a a a todo!",
				State: false,
			},
			{
				Id:    4,
				Title: "Dont dont care",
				Note:  "!!!!!!!!!!!!!!!!!!!!!!!!!!!!",
				State: false,
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
		fmt.Println("State: " + strconv.FormatBool(list.TodoList[i].State))
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

	return tdList
}
func printItem(item Todo, w http.ResponseWriter) {
	fmt.Fprintln(w, "---------------------------------------")
	fmt.Fprintln(w, "ID: "+strconv.Itoa(item.Id))
	fmt.Fprintln(w, "Title: "+item.Title)
	fmt.Fprintln(w, "Note: "+item.Note)
	fmt.Fprintln(w, "State: "+strconv.FormatBool(item.State))
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

func generateNewID(list TodoList) int {

	max := 0

	//generation of id(last id at the json file+1)
	for _, item := range list.TodoList {
		if item.Id > max {
			max = item.Id
		}
	}
	return (max + 1)
}

func addToList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: Add new item")

	list := readList()

	if r.Method == "GET" {
		t, _ := template.ParseFiles("form.html")
		t.Execute(w, nil)
		return
	} else {
		r.ParseForm()

		id := generateNewID(list)
		title := r.Form["title"]
		body := r.Form["body"]
		fmt.Fprintln(w, "id: ", id)
		fmt.Fprintln(w, "title: ", title[0])
		fmt.Fprintln(w, "body: ", body[0])

		fmt.Println("id: ", id)
		fmt.Println("title: ", title[0])
		fmt.Println("body: ", body[0])

		item := Todo{
			Id:    id,
			Title: title[0],
			Note:  body[0],
			State: false,
		}

		list.TodoList = append(list.TodoList, item)
		file, _ := json.MarshalIndent(list, "", "")

		fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		fmt.Println("Finished writing")

	}

}

func deleteFromList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: Delete item")

	list := readList()

	fmt.Println("Length: ", len(list.TodoList))

	if r.Method == "GET" {
		t, _ := template.ParseFiles("formDel.html")
		t.Execute(w, list)

		return
	} else {
		r.ParseForm()
		idstr := r.Form["id"]
		ind, _ := strconv.Atoi(idstr[0])
		fmt.Println("id: ", ind)

		for index := range list.TodoList {
			fmt.Println("Looping")

			if index == ind {
				fmt.Println("Deleting: ")

				list.TodoList = append(list.TodoList[:index], list.TodoList[index+1:]...)
			}
		}
		printList(list, w)
		file, _ := json.MarshalIndent(list, "", "")
		fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		fmt.Println("Finished writing")
	}
}
func updateItem(w http.ResponseWriter, r *http.Request) {
	list := readList()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("ID is: ", id)
	}
	var index int
	var found bool = false
	for i := 0; i < len(list.TodoList); i++ {
		fmt.Println(id)
		if id == list.TodoList[i].Id {
			index = i
			found = true
			break
		}
	}

	if r.Method == "GET" {

		if found {
			t, _ := template.ParseFiles("formUpdate.html")
			t.Execute(w, list.TodoList[index])
		} else {
			t, _ := template.ParseFiles("formUpdateNotFound.html")
			t.Execute(w, nil)
		}

		return
	} else {
		r.ParseForm()

		title := r.Form["Title"]
		body := r.Form["Body"]
		state := r.Form["State"]
		s, err := strconv.ParseBool(state[0])

		fmt.Println("Form: ", r.Form)

		list.TodoList[index].Title = title[0]
		list.TodoList[index].Note = body[0]
		if err == nil {
			list.TodoList[index].State = s
		}

		printItem(list.TodoList[index], w)

		file, _ := json.MarshalIndent(list, "", "")
		fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		fmt.Println("Finished writing")
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/showAll", getFullList)
	myRouter.HandleFunc("/show/{id}", getOneItem)
	myRouter.HandleFunc("/add", addToList)
	myRouter.HandleFunc("/del", deleteFromList)
	myRouter.HandleFunc("/update/{id}", updateItem)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func init() {
	if _, err := os.Stat("storage/todoList.json"); os.IsNotExist(err) {
		initStore()
	}
}
func main() {

	readList()

	handleRequests()

}
