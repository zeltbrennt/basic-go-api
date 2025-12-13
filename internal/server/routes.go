package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/zeltbrennt/go-api/docs"
)

func (ts *TaskServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", ts.getAllTasksHandler)
	mux.HandleFunc("POST /tasks", ts.createTaskHandler)
	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json")))
	return mux
}
