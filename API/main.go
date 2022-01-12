package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// task structs
type Task struct {
	ID          string `json: "id"`
	TaskName    string `json: "taskname"`
	Description string `json: "description"`
	CreatedOn   string `json: "createdon"`
	Priority    string `json: "priority"`
	Creator     *User  `json: "creator"`
}

type User struct {
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
}

var tasks []Task

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, task := range tasks {
		if task.ID == params["id"] {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	json.NewEncoder(w).Encode("No task of such ID!")
}
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	task.ID = strconv.Itoa(rand.Intn(100))
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
}
func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, task := range tasks {
		if task.ID == params["id"] {
			tasks = append(tasks[:i], tasks[i+1:]...)
			var task Task
			_ = json.NewDecoder(r.Body).Decode(&task)
			task.ID = params["id"]
			tasks = append(tasks, task)
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	json.NewEncoder(w).Encode("No such task to update!")
}
func removeTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, task := range tasks {
		if task.ID == params["id"] {
			tasks = append(tasks[:i], tasks[i+1:]...)
			json.NewEncoder(w).Encode("Task Deleted.")
			return
		}
	}
	json.NewEncoder(w).Encode("No such task to delete!")
}

func main() {

	// init router
	r := mux.NewRouter()

	//fake tasks
	tasks = append(tasks, Task{ID: "1",
		TaskName:    "Learn GO",
		Description: "Learn how to create a REST API using GO!",
		CreatedOn:   "13/12/2021",
		Priority:    "HIGH!",
		Creator: &User{FirstName: "Ryan",
			LastName: "Low"}})

	// route handler
	r.HandleFunc("/api/tasks", getTasks).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/api/tasks", createTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", removeTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
