package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var todoList []todo //Save todos
var counter int     //Keep track of todo number

func viewTodoListEndpoint(w http.ResponseWriter, r *http.Request) {
	for _, todo := range todoList {
		var temp []byte
		temp = []byte(fmt.Sprintf("ID: %d Title: %s, Description: %s\n", todo.ID, todo.Title, todo.Description))
		w.Write(temp)
	}
}

func addTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	var newTodo todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	counter++
	newTodo.ID = counter

	todoList = append(todoList, newTodo)

	var temp []byte
	temp = []byte(fmt.Sprintf("TODO with id %d added succesfully\n", newTodo.ID))
	w.Write(temp)
	log.Println(string(temp))
}

func deleteTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	var id todo
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for i, todo := range todoList {
		if todo.ID == id.ID {
			todoList = append(todoList[:i], todoList[i+1:]...)
			log.Printf("TODO with id %d deleted succesfully\n", todo.ID)
			break
		}
	}
	for _, todo := range todoList {
		var temp []byte
		temp = []byte(fmt.Sprintf("ID: %d Title: %s, Description: %s\n", todo.ID, todo.Title, todo.Description))
		w.Write(temp)
	}
}

func main() {
	var router = mux.NewRouter()

	router.HandleFunc("/todo", viewTodoListEndpoint).Methods("GET")
	router.HandleFunc("/todo/new", addTodoEndpoint).Methods("POST")
	router.HandleFunc("/todo/delete", deleteTodoEndpoint).Methods("POST")

	address := ":8000"
	log.Printf("Starting server on localhost%s\n\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Println(err)
	}
}
