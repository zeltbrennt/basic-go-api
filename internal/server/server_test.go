package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeltbrennt/go-api/internal/models"
	"github.com/zeltbrennt/go-api/internal/server"
	"github.com/zeltbrennt/go-api/internal/store"
)

var mockStore = store.NewMockStore()

func newTestServer() http.Handler {
	server := server.NewTaskService(mockStore)
	return server.Routes()
}

func doRequest(ts http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	ts.ServeHTTP(w, req)
	return w
}

func TestGetAllTasks(t *testing.T) {
	t.Log("testing all Tasks")
	ts := newTestServer()
	for i := range 3 {
		task := models.Task{
			ID:    i,
			Title: fmt.Sprintf("Task #%d", i),
		}
		_, err := mockStore.CreateTask(task)
		if err != nil {
			t.Fatal("error while creating task")
		}
	}

	res := doRequest(ts, "GET", "/tasks", nil)

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
