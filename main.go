package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

var Todos []Todo = []Todo{
	{Id: generateId(), Title: "Buy groceries", Description: "Milk, bread, eggs", IsCompleted: false},
	{Id: generateId() + 1, Title: "Finish report", Description: "Due tomorrow", IsCompleted: false},
}

func generateId() int {
	return int(time.Now().UnixNano())
}

func main() {
	http.HandleFunc("/", HealthCheck)
	http.HandleFunc("/all-todos", AllTodos)
	http.HandleFunc("/add-todo", AddTodo)
	http.HandleFunc("/get-todo/", GetTodoById)
	http.HandleFunc("/delete-todo/", DeleteTodo)
	http.HandleFunc("/update-todo/", UpdateTodo)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// HealthCheck checks the health of the server
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATH::", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("Server health is okk")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AllTodos returns all todos
func AllTodos(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATH::", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/all-todos") {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(Todos)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AddTodo adds a new todo
func AddTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATH::", r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/add-todo") {
		http.NotFound(w, r)
		return
	}
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	newTodo.Id = generateId()
	Todos = append(Todos, newTodo)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode("Todo Added Successfully")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetTodoById returns a todo by its ID
func GetTodoById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATH::", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/get-todo/") {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/get-todo/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	todoId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	for index, todo := range Todos {
		if todo.Id == todoId {
			todo := Todos[index]
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(todo)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	err = json.NewEncoder(w).Encode("No Todo found with this id")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DeleteTodo deletes a todo by its ID
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATH::", r.URL.Path)
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/delete-todo/") {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/delete-todo/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	todoId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	for index, todo := range Todos {
		if todo.Id == todoId {
			Todos = append(Todos[:index], Todos[index+1:]...)
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode("Todo deleted successfully")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	err = json.NewEncoder(w).Encode("No Todo found with this id")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// UpdateTodo data -> is_completed true | false
type UpdateTodoRequest struct {
	IsCompleted bool `json:"is_completed"`
}

// UpdateTodo updates a todo by its ID
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PATH::", r.URL.Path)
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.URL.Path, "/update-todo/") {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/update-todo/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}
	todoId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	var updateData UpdateTodoRequest
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Please provide enough data to update todo", http.StatusBadRequest)
		return
	}
	for index, todo := range Todos {
		if todo.Id == todoId {
			Todos[index].IsCompleted = updateData.IsCompleted
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode("Todo Updated successfully")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	err = json.NewEncoder(w).Encode("No Todo found with this id")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
