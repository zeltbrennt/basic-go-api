// Package server
//
// Describes the server with its routes and handers
package server

import "github.com/zeltbrennt/go-api/internal/store"

type TaskServer struct {
	store store.Store
}

func NewTaskServer(s store.Store) *TaskServer {
	return &TaskServer{
		store: s,
	}
}
