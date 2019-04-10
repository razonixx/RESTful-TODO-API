package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var todoList []todo //Save todos
var counter int     //Keep track of todo number

func indexEndpoint(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func viewTodoListEndpoint(w http.ResponseWriter, r *http.Request) {
	for _, todo := range todoList {
		fmt.Fprintf(w, "ID: %d Title: %s, Description: %s\n", todo.ID, todo.Title, todo.Description)
	}
}

func viewTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	targetTodo := mux.Vars(r)
	id, err := strconv.Atoi(targetTodo["todoId"])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	for _, todo := range todoList {
		if todo.ID == id {
			fmt.Fprintf(w, "ID: %d Title: %s, Description: %s\n", todo.ID, todo.Title, todo.Description)
			return
		}
	}
	fmt.Fprintf(w, "No TODO with ID %d was found\n", id)
}

func addTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	var newTodo todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	counter++
	newTodo.ID = counter

	todoList = append(todoList, newTodo)

	fmt.Fprintf(w, "TODO with ID %d added succesfully!\n", newTodo.ID)
	log.Printf("TODO with id %d added succesfully\n", newTodo.ID)
}

func deleteTodoEndpoint(w http.ResponseWriter, r *http.Request) {
	targetTodo := mux.Vars(r)
	id, err := strconv.Atoi(targetTodo["todoId"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for i, todo := range todoList {
		if todo.ID == id {
			todoList = append(todoList[:i], todoList[i+1:]...)
			fmt.Fprintf(w, "TODO with ID %d deleted succesfully!\n", i)
			log.Printf("TODO with id %d deleted succesfully\n", todo.ID)
			return
		}
	}
	fmt.Fprintf(w, "No TODO with ID %d was found\n", id)
}

func main() {
	var router = mux.NewRouter()

	var newTodo todo
	newTodo.ID = counter
	newTodo.Title = "Hola"
	newTodo.Description = "Desc"

	todoList = append(todoList, newTodo)

	router.HandleFunc("/", indexEndpoint).Methods("GET")
	router.HandleFunc("/todo", viewTodoListEndpoint).Methods("GET")
	router.HandleFunc("/todo/{todoId}", viewTodoEndpoint).Methods("GET")
	router.HandleFunc("/todo/new", addTodoEndpoint).Methods("POST").Headers("Content-Type", "application/json")
	router.HandleFunc("/todo/delete/{todoId}", deleteTodoEndpoint).Methods("POST")

	address := ":8000"
	log.Printf("Starting server on localhost%s\n\n", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Println(err)
	}
}
