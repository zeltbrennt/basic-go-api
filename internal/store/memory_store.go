package store

import (
	"context"
	"sync"

	"github.com/zeltbrennt/go-api/internal/models"
)

type memoryStore struct {
	mu    sync.Mutex
	tasks map[int]models.Task
	next  int
}

func (m *memoryStore) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]models.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (m *memoryStore) CreateTask(ctx context.Context, t models.Task) (models.Task, error) {
	select {
	case <-ctx.Done():
		return models.Task{}, ctx.Err()
	default:
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	t.ID = m.next
	m.tasks[t.ID] = t
	m.next++
	return t, nil
}

func NewMockStore() *memoryStore {
	return &memoryStore{
		tasks: make(map[int]models.Task),
		next:  1,
	}
}
