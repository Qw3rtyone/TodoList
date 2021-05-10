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
			{
				Id:    5,
				Title: "Number 5",
				Note:  "Now number 5",
				State: true,
			},
			{
				Id:    6,
				Title: "Number 6",
				Note:  "Not a things",
				State: true,
			},
			{
				Id:    7,
				Title: "Oh god how many more",
				Note:  "Get comfortableGet comfortableGet comfortableGet comfortableGet comfortableGet comfortable",
				State: false,
			},
			{
				Id:    8,
				Title: "T",
				Note:  "T",
				State: false,
			},
			{
				Id:    9,
				Title: "On this earth",
				Note:  "There lay a child in a grave",
				State: false,
			},
			{

				Id:    10,
				Title: "Friday: Milk",
				Note:  "Buy milk for friday",
				State: false,
			},
			{
				Id:    11,
				Title: "Dog",
				Note:  "Buy a dog",
				State: true,
			},
			{
				Id:    12,
				Title: "Cat",
				Note:  "Sell the cat",
				State: false,
			},
		},
	}

	file, _ := json.MarshalIndent(data, "", "")
	fmt.Println("writing")
	ioutil.WriteFile("storage/todoList.json", file, 0644)
	fmt.Println("Finished writing")
}

func readList() TodoList {
	data, err := os.Open("storage/todoList.json")
	checkerr(err)
	fmt.Println("Opened file!!")

	defer data.Close()

	byteVal, _ := ioutil.ReadAll(data)

	var tdList TodoList

	json.Unmarshal(byteVal, &tdList)

	return tdList
}

func homePage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/homepage.html")
	t.Execute(w, nil)

	fmt.Println("End of page")
}
func getFullList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: full list")

	t, _ := template.ParseFiles("templates/showAllItems.html")
	t.Execute(w, readList())
}
func getOneItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: single item")
	vars := mux.Vars(r)
	key := vars["id"]

	for _, item := range readList().TodoList {
		k, err := strconv.Atoi(key)
		fmt.Println(k)
		if (err == nil) && (k == item.Id) {
			t, _ := template.ParseFiles("templates/singleItem.html")
			t.Execute(w, item)
			return
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
		t, _ := template.ParseFiles("templates/form.html")
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
		t, _ := template.ParseFiles("templates/formDel.html")
		t.Execute(w, list)

		return
	}
	if r.Method == "POST" {
		r.ParseForm()

		idstr := r.Form["idbutton"]
		ind, _ := strconv.Atoi(idstr[0])

		fmt.Println("id: ", ind)

		for index, item := range list.TodoList {
			fmt.Println("Looping")

			if item.Id == ind {
				fmt.Println("Deleting: ")

				list.TodoList = append(list.TodoList[:index], list.TodoList[index+1:]...)
			}
		}

		t, _ := template.ParseFiles("templates/formDel.html")
		t.Execute(w, list)

		file, _ := json.MarshalIndent(list, "", "")
		fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		fmt.Println("Finished writing")

		return
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
			t, _ := template.ParseFiles("templates/formUpdate.html")
			t.Execute(w, list.TodoList[index])
		} else {
			t, _ := template.ParseFiles("templates/formUpdateNotFound.html")
			t.Execute(w, nil)
		}

		return
	} else {
		r.ParseForm()

		title := r.Form["Title"]
		body := r.Form["Note"]
		state := r.Form["State"]
		s, err := strconv.ParseBool(state[0])

		fmt.Println("Form: ", r.Form)

		list.TodoList[index].Title = title[0]
		list.TodoList[index].Note = body[0]
		if err == nil {
			list.TodoList[index].State = s
		}

		t, _ := template.ParseFiles("templates/singleItem.html")
		t.Execute(w, list.TodoList[index])

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

	myRouter.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css/"))))
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
