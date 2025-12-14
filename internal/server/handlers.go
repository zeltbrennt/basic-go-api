package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zeltbrennt/go-api/internal/models"
)

var validate = validator.New()

// GetAllTasks godoc
//
//	@Summary		Gets all tasks
//	@Description	Gets all tasks
//	@Tags			tasks
//	@Produce		json
//	@Success		200	{array}	models.Task
//	@Failure		500
//	@Router			/tasks [get]
func (ts *TaskService) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ts.logger.Info("get all tasks", "path", r.URL.Path)
	tasks, err := ts.store.GetAllTasks(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "error while decoding json", http.StatusInternalServerError)
		return
	}
}

// CreateTask godoc
//
//	@Summary		Create a new Task
//	@Description	Create a task and returns it
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		models.Task	true	"Some Task"
//	@Success		201		{object}	models.Task
//	@Failure		400
//	@Router			/tasks [post]
func (ts *TaskService) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		ts.logger.Warn("invalid json", "path", r.URL)
		return
	}
	// validate json
	if err := validate.Struct(t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		ts.logger.Warn("invalid json", "path", r.URL)
		return
	}
	newTask, err := ts.store.CreateTask(ctx, t)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			http.Error(w, "request canceled", http.StatusRequestTimeout)
			ts.logger.Warn("request canceled", "path", r.URL)
			return
		}
		ts.logger.Error("error while saving task", "error", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(newTask)
	if err != nil {
		ts.logger.Error("error while decoding json", "error", err.Error())
		http.Error(w, "error while decoding json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
