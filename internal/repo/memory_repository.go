package repo

import (
	"context"
	"sync"
	"time"
	"github.com/google/uuid"
)

// MemoryTaskRepository реализует TaskRepository в памяти
// Используется для демонстрации без необходимости в базе данных
// Потокобезопасен благодаря использованию RWMutex
type MemoryTaskRepository struct {
	tasks []*Task
	mutex sync.RWMutex
}

// NewMemoryTaskRepository создает новый экземпляр in-memory репозитория
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make([]*Task, 0),
	}
}

// Create добавляет новую задачу в память
// Генерирует UUID и устанавливает время создания
func (r *MemoryTaskRepository) Create(ctx context.Context, task *Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	r.tasks = append(r.tasks, task)
	return nil
}

// GetAll возвращает все задачи из памяти
// Создает копии задач для предотвращения изменения оригинальных данных
func (r *MemoryTaskRepository) GetAll(ctx context.Context) ([]*Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// Создаем копию слайса для безопасности
	result := make([]*Task, len(r.tasks))
	for i, task := range r.tasks {
		// Создаем копию задачи
		taskCopy := *task
		result[i] = &taskCopy
	}
	return result, nil
}

func (r *MemoryTaskRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, task := range r.tasks {
		if task.ID == id {
			// Создаем копию задачи
			taskCopy := *task
			return &taskCopy, nil
		}
	}
	return nil, nil
}

func (r *MemoryTaskRepository) Update(ctx context.Context, task *Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = task
			return nil
		}
	}
	return nil
}

func (r *MemoryTaskRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return nil
}
