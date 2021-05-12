package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
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
		t.Errorf("The readList function did not return the list equal to provided")
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

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
