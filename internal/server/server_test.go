package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeltbrennt/go-api/internal/models"
	"github.com/zeltbrennt/go-api/internal/server"
	"github.com/zeltbrennt/go-api/internal/store"
)

func newTestServer() (http.Handler, store.Storer) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	store := store.NewMockStore()
	server := server.NewTaskService(store, logger)
	return server.Routes(), store
}

func doRequest(ts http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	return w
}

func TestGetAllTasks(t *testing.T) {
	ctx := context.Background()
	ts, mockStore := newTestServer()
	for i := range 3 {
		task := models.Task{
			ID:    i,
			Title: fmt.Sprintf("Task #%d", i),
		}
		_, err := mockStore.CreateTask(ctx, task)
		if err != nil {
			t.Fatal("error while creating task")
		}
	}

	res := doRequest(ts, "GET", "/api/v1/tasks", nil)

	if res.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, res.Code)
	}

	var tasks []models.Task
	err := json.NewDecoder(res.Body).Decode(&tasks)
	if err != nil {
		t.Fatal("error while decoding json")
	}
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}
