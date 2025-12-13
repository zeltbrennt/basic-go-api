package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeltbrennt/go-api/internal/models"
)

var validate = validator.New()

// GetAllTasks godoc
// @Summary Gets all tasks
// @Description Gets all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500
// @Router /tasks [get]
func (ts *TaskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("get all tasks %s\n", r.URL.Path)
	tasks, _ := ts.store.GetAllTasks()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "error while decoding json", http.StatusInternalServerError)
		return
	}
}

// CreateTask godoc
// @Summary Create a new Task
// @Description Create a task and returns it
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Some Task"
// @Success 201 {object} models.Task
// @Failure 400
// @Router /tasks [post]
func (ts *TaskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	// validate json
	if err := validate.Struct(t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	newTask, err := ts.store.CreateTask(t)
	if err != nil {
		http.Error(w, "error while saving", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newTask)
	if err != nil {
		http.Error(w, "error while decoding json", http.StatusInternalServerError)
	}
}
