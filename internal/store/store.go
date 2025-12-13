// Package store
//
// Describes the interface the storage needs to fullfill
package store

import "github.com/zeltbrennt/go-api/internal/models"

type Store interface {
	GetAllTasks() ([]models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
}
