package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGenerateDefaultList(t *testing.T) {
	got := generateDefaultList()
	want := TodoList{
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

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Generated struct not equal to default")
	}
}

func TestInitStore(t *testing.T) {
	e := os.Remove("storage/todoList.json")

	if e != nil {
		t.Fatal(e)
	}

	got := initStore()

	if got != nil {
		t.Errorf("Path doesn't exist, json file not created")
	}
	if _, err := os.Stat("storage/todoList.json"); os.IsNotExist(err) {
		t.Errorf("file not found")
	}

}

func TestReadList(t *testing.T) {

	got := readList()

	data, err := os.Open("storage/todoList.json")
	defer data.Close()

	if err != nil {
		t.Errorf("Couldn't read comparison list")
	}

	byteVal, _ := ioutil.ReadAll(data)

	var want TodoList
	json.Unmarshal(byteVal, &want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("The readList function did not return the list equal to List file")
	}
}

type testStruct struct {
	todolist TodoList
	expected int
}

var idTests = []testStruct{
	{
		todolist: TodoList{
			TodoList: []Todo{
				{
					Id:    0,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
			},
		},
		expected: 1,
	},
	{
		todolist: TodoList{
			TodoList: []Todo{
				{
					Id:    10,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
			},
		},
		expected: 11,
	},
	{
		todolist: TodoList{
			TodoList: []Todo{
				{
					Id:    0,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    1,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    2,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    3,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    4,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
			},
		},
		expected: 5,
	},
	{
		todolist: TodoList{
			TodoList: []Todo{
				{
					Id:    100,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    11231,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    9999,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    98182,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
				{
					Id:    4,
					Title: "Default",
					Note:  "Default Note",
					State: false,
				},
			},
		},
		expected: 98183,
	},
}

func TestGenerateNewID(t *testing.T) {
	t.Log("Generate ID")
	for _, test := range idTests {
		got := generateNewID(test.todolist)
		if got != test.expected {
			t.Errorf("Expected Id is not generated. Expected: %d, recieved: %d", test.expected, got)
		}
	}

}

func TestGetFullList(t *testing.T) {
	req, err := http.NewRequest("GET", "/showAll", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getFullList)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	a, _ := template.ParseFiles("templates/showAllItems.html")
	var tpl bytes.Buffer
	a.Execute(&tpl, readList())

	// Check the response body is what we expect.
	expected := tpl.String()
	//t.Errorf("Print: want %v  got %v", expected, rr.Body.String())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestGetOneItem(t *testing.T) {

	testRequests := []string{"0", "1", "3"}

	list := readList()

	for _, id := range testRequests {
		req, err := http.NewRequest("GET", "/show/"+id, nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/show/{id}", getOneItem)
		router.ServeHTTP(rr, req)

		//check status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		listId, err := strconv.Atoi(id)
		if err != nil {
			t.Errorf("Couldn't convert key to int")
		}
		var listInd int
		for index, item := range list.TodoList {
			if item.Id == listId {
				listInd = index
				break
			}
		}

		a, _ := template.ParseFiles("templates/singleItem.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, list.TodoList[listInd])

		got := rr.Body.String()
		want := tpl.String()

		if got != want {
			t.Errorf("handler returned unexpected body. got %v want %v", got, want)
		}
	}
}

func TestGetOneItemNotFound(t *testing.T) {

	testRequests := []string{"9999", "-1"}

	for _, id := range testRequests {
		req, err := http.NewRequest("GET", "/show/"+id, nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/show/{id}", getOneItem)
		router.ServeHTTP(rr, req)

		//check status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		a, _ := template.ParseFiles("templates/formNotFound.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, nil)

		got := rr.Body.String()
		want := tpl.String()

		if got != want {
			t.Errorf("handler returned unexpected body. got %v want %v", got, want)
		}
	}
}

func TestHomepage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homePage)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	a, _ := template.ParseFiles("templates/homepage.html")
	var tpl bytes.Buffer
	a.Execute(&tpl, readList())

	// Check the response body is what we expect.
	want := tpl.String()
	//t.Errorf("Print: want %v  got %v", expected, rr.Body.String())
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}

func TestAddToListGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/add", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addToList)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	a, _ := template.ParseFiles("templates/form.html")
	var tpl bytes.Buffer
	a.Execute(&tpl, readList())

	// Check the response body is what we expect.
	want := tpl.String()
	//t.Errorf("Print: want %v  got %v", expected, rr.Body.String())
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), want)
	}
}

func TestAddToListPost(t *testing.T) {
	testFormStruct := []struct {
		testTitle string
		testBody  string
		expected  Todo
	}{
		{
			"Test1",
			"Test1 body",
			Todo{
				4,
				"Test1",
				"Test1 body",
				false,
			},
		},
		{
			"Test2",
			"Test2 body. This is a test body",
			Todo{
				5,
				"Test2",
				"Test2 body. This is a test body",
				false,
			},
		},
		{
			"Test 3",
			"Test 3 body. This is a test body for the third test.",
			Todo{
				6,
				"Test 3",
				"Test 3 body. This is a test body for the third test.",
				false,
			},
		},
	}

	e := os.Remove("storage/todoList.json")
	if e != nil {
		t.Fatal(e)
	}
	err := initStore()
	if err != nil {
		t.Fatal(e)
	}

	for _, test := range testFormStruct {
		form := url.Values{}
		form.Add("title", test.testTitle)
		form.Add("body", test.testBody)

		req, err := http.NewRequest("POST", "/add", strings.NewReader(form.Encode()))
		req.Form = form
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/add", addToList)
		router.ServeHTTP(rr, req)

		//check status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		a, _ := template.ParseFiles("templates/singleItem.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, test.expected)

		got := rr.Body.String()
		want := tpl.String()

		if got != want {
			t.Errorf("handler returned unexpected body. got %v want %v", got, want)
		}
	}

}

func TestDeleteFromListGet(t *testing.T) {
	e := os.Remove("storage/todoList.json")
	if e != nil {
		t.Fatal(e)
	}
	err := initStore()
	if err != nil {
		t.Fatal(e)
	}

	req, err := http.NewRequest("GET", "/del", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteFromList)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	a, _ := template.ParseFiles("templates/formDel.html")
	var tpl bytes.Buffer
	a.Execute(&tpl, readList())

	// Check the response body is what we expect.
	got := rr.Body.String()
	want := tpl.String()
	if got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}

}

func TestDeleteFromListPost(t *testing.T) {
	e := os.Remove("storage/todoList.json")
	if e != nil {
		t.Fatal(e)
	}
	err := initStore()
	if err != nil {
		t.Fatal(e)
	}

	testDelStruct := []struct {
		delID    string
		expected TodoList
	}{
		{
			delID: "0",
			expected: TodoList{
				TodoList: []Todo{
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
			},
		},
		{
			delID: "2",
			expected: TodoList{
				TodoList: []Todo{
					{
						Id:    1,
						Title: "Add New Todo",
						Note:  "Add a third todo!",
						State: true,
					},
					{
						Id:    3,
						Title: "Add Next",
						Note:  "Add a a a a todo!",
						State: false,
					},
				},
			},
		},
		{
			delID: "3",
			expected: TodoList{
				TodoList: []Todo{
					{
						Id:    1,
						Title: "Add New Todo",
						Note:  "Add a third todo!",
						State: true,
					},
				},
			},
		},
	}

	for _, testcase := range testDelStruct {

		form := url.Values{}
		form.Add("idbutton", testcase.delID)
		req, err := http.NewRequest("POST", "/del", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Form = form

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/del", deleteFromList)
		router.ServeHTTP(rr, req)

		a, _ := template.ParseFiles("templates/formDel.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, testcase.expected)

		// Check the response body is what we expect.
		got := rr.Body.String()
		want := tpl.String()
		if got != want {
			t.Errorf("handler returned unexpected body: got %v want %v", got, want)
		}
	}

}

func TestUpdateItemGet(t *testing.T) {
	//clean up the damage caused by the previous test and reset to default state
	e := os.Remove("storage/todoList.json")
	if e != nil {
		t.Fatal(e)
	}
	err := initStore()
	if err != nil {
		t.Fatal(e)

	}

	testUpdateGetStruct := []struct {
		id       string
		expected Todo
	}{
		{
			id: "0",
			expected: Todo{
				Id:    0,
				Title: "Default",
				Note:  "Default Note",
				State: false,
			},
		},
		{
			id: "3",
			expected: Todo{
				Id:    3,
				Title: "Add Next",
				Note:  "Add a a a a todo!",
				State: false,
			},
		},
		{
			id: "1",
			expected: Todo{
				Id:    1,
				Title: "Add New Todo",
				Note:  "Add a third todo!",
				State: true,
			},
		},
	}
	for _, testcase := range testUpdateGetStruct {

		req, err := http.NewRequest("GET", "/update/"+testcase.id, nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/update/{id}", updateItem)
		router.ServeHTTP(rr, req)

		//check status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		a, _ := template.ParseFiles("templates/formUpdate.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, testcase.expected)

		got := rr.Body.String()
		want := tpl.String()

		if got != want {
			t.Errorf("handler returned unexpected body: got %v want %v", got, want)
		}

	}

}

func TestUpdateItemNotFound(t *testing.T) {
	//clean up the damage caused by the previous test and reset to default state
	e := os.Remove("storage/todoList.json")
	if e != nil {
		t.Fatal(e)
	}
	err := initStore()
	if err != nil {
		t.Fatal(e)

	}

	testRequests := []string{"9999", "-1", "112312312"}

	for _, id := range testRequests {
		req, err := http.NewRequest("GET", "/update/"+id, nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/update/{id}", updateItem)
		router.ServeHTTP(rr, req)

		//check status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		a, _ := template.ParseFiles("templates/formNotFound.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, nil)

		got := rr.Body.String()
		want := tpl.String()

		if got != want {
			t.Errorf("handler returned unexpected body. got %v want %v", got, want)
		}
	}

}

func TestUpdateItemPost(t *testing.T) {
	//clean up the damage caused by the previous test and reset to default state
	e := os.Remove("storage/todoList.json")
	if e != nil {
		t.Fatal(e)
	}
	err := initStore()
	if err != nil {
		t.Fatal(e)

	}
	testUpdateStruct := []struct {
		updateID    string
		updateTitle string
		updateBody  string
		updateState string
		expected    Todo
	}{
		{
			updateID:    "0",
			updateTitle: "New Title for id 0",
			updateBody:  "New body for id 0",
			updateState: "true",
			expected: Todo{
				Id:    0,
				Title: "New Title for id 0",
				Note:  "New body for id 0",
				State: true,
			},
		},
		{
			updateID:    "0",
			updateTitle: "New Title for id 0",
			updateBody:  "New body for id 0",
			updateState: "false",
			expected: Todo{
				Id:    0,
				Title: "New Title for id 0",
				Note:  "New body for id 0",
				State: false,
			},
		},
		{
			updateID:    "2",
			updateTitle: "I want to change thissss",
			updateBody:  "NP",
			updateState: "false",
			expected: Todo{
				Id:    2,
				Title: "I want to change thissss",
				Note:  "NP",
				State: false,
			},
		},
	}

	for _, testcase := range testUpdateStruct {
		form := url.Values{}
		form.Add("Title", testcase.updateTitle)
		form.Add("Note", testcase.updateBody)
		form.Add("State", testcase.updateState)
		req, err := http.NewRequest("POST", "/update/"+testcase.updateID, strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Form = form

		rr := httptest.NewRecorder()
		//create a new router to make sure the variables are passed in properly
		router := mux.NewRouter()
		router.HandleFunc("/update/{id}", updateItem)
		router.ServeHTTP(rr, req)

		//check status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		a, _ := template.ParseFiles("templates/singleItem.html")
		var tpl bytes.Buffer
		a.Execute(&tpl, testcase.expected)

		got := rr.Body.String()
		want := tpl.String()

		if got != want {
			t.Errorf("handler returned unexpected body. got %v want %v", got, want)
		}
	}

	//t.Errorf("Todo")
}
