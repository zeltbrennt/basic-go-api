// Package store
//
// Describes the interface the storage needs to fullfill
package store

import (
	"context"

	"github.com/zeltbrennt/go-api/internal/models"
)

type Storer interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	CreateTask(ctx context.Context, task models.Task) (models.Task, error)
}
