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

func generateDefaultList() TodoList {
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
				State: true,
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
		},
	}
	return data
}
func initStore() error {
	if _, err := os.Stat("storage/todoList.json"); os.IsNotExist(err) {

		data := generateDefaultList()
		file, err := json.MarshalIndent(data, "", "")

		if err != nil {
			return err
		}

		//fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		//fmt.Println("Finished writing")
		return nil
	}

	return nil
}

func readList() TodoList {
	data, err := os.Open("storage/todoList.json")
	checkerr(err)
	//fmt.Println("Opened file!!")

	defer data.Close()

	byteVal, _ := ioutil.ReadAll(data)

	var tdList TodoList

	json.Unmarshal(byteVal, &tdList)

	return tdList
}

func homePage(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("templates/homepage.html")
	t.Execute(w, nil)

	//fmt.Println("End of page")
}
func getFullList(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint: full list")

	t, _ := template.ParseFiles("templates/showAllItems.html")
	t.Execute(w, readList())

}
func getOneItem(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint: single item")
	vars := mux.Vars(r)
	key := vars["id"]
	k, err := strconv.Atoi(key)
	for _, item := range readList().TodoList {
		//fmt.Println(k)
		if (err == nil) && (k == item.Id) {
			t, _ := template.ParseFiles("templates/singleItem.html")
			t.Execute(w, item)
			return
		}
	}
	t, _ := template.ParseFiles("templates/formNotFound.html")
	t.Execute(w, nil)
	return
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
	//fmt.Println("Endpoint: Add new item")

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

		//fmt.Println("id: ", id)
		//fmt.Println("title: ", title[0])
		//fmt.Println("body: ", body[0])

		item := Todo{
			Id:    id,
			Title: title[0],
			Note:  body[0],
			State: false,
		}

		list.TodoList = append(list.TodoList, item)
		file, _ := json.MarshalIndent(list, "", "")

		t, _ := template.ParseFiles("templates/singleItem.html")
		t.Execute(w, item)

		//fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		//fmt.Println("Finished writing")

	}

}

func deleteFromList(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint: Delete item")

	list := readList()

	//fmt.Println("Length: ", len(list.TodoList))

	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/formDel.html")
		t.Execute(w, list)

		return
	}
	if r.Method == "POST" {
		r.ParseForm()

		idstr := r.Form["idbutton"]
		ind, _ := strconv.Atoi(idstr[0])

		//fmt.Println("id: ", ind)

		for index, item := range list.TodoList {
			//fmt.Println("Looping")

			if item.Id == ind {
				//fmt.Println("Deleting: ")

				list.TodoList = append(list.TodoList[:index], list.TodoList[index+1:]...)
			}
		}

		t, _ := template.ParseFiles("templates/formDel.html")
		t.Execute(w, list)

		file, _ := json.MarshalIndent(list, "", "")
		//fmt.Println("writing")
		ioutil.WriteFile("storage/todoList.json", file, 0644)
		//fmt.Println("Finished writing")

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
		//fmt.Println(id)
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
			t, _ := template.ParseFiles("templates/formNotFound.html")
			t.Execute(w, nil)
		}
		return

	} else {
		r.ParseForm()

		title := r.Form["Title"]
		body := r.Form["Note"]
		state := r.Form["State"]
		s, err := strconv.ParseBool(state[0])

		//fmt.Println("Form: ", r.Form)

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

	err := initStore()
	checkerr(err)
}
func main() {

	readList()

	handleRequests()

}
