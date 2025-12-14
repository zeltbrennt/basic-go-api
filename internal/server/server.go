// Package server
//
// Describes the server with its routes and handers
package server

import (
	"log/slog"

	"github.com/zeltbrennt/go-api/internal/store"
)

type TaskService struct {
	store  store.Storer
	logger *slog.Logger
}

func NewTaskService(store store.Storer, logger *slog.Logger) *TaskService {
	return &TaskService{
		store:  store,
		logger: logger,
	}
}
