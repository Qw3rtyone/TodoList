package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
		log.Fatal(e)
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
