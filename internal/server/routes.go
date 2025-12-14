package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zeltbrennt/go-api/docs"
	"github.com/zeltbrennt/go-api/internal/middleware"
)

func (ts *TaskService) Routes() http.Handler {
	rootMux := http.NewServeMux()
	rootMux.Handle("/api/v1/", handleAPIv1(ts))
	rootMux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json")))
	return rootMux
}

// Routes to /api/v1
func handleAPIv1(ts *TaskService) http.Handler {
	api := http.NewServeMux()
	api.HandleFunc("GET /tasks", ts.getAllTasksHandler)
	api.HandleFunc("POST /tasks", ts.createTaskHandler)
	return middleware.LoggingMiddleware(ts.logger)(http.StripPrefix("/api/v1", api))
}
