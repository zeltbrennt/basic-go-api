package store

import (
	"sync"

	"github.com/zeltbrennt/go-api/internal/models"
)

type memoryStore struct {
	mu    sync.Mutex
	tasks map[int]models.Task
	next  int
}

func (m *memoryStore) GetAllTasks() ([]models.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var result []models.Task
	for _, t := range m.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (m *memoryStore) CreateTask(t models.Task) (models.Task, error) {
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
