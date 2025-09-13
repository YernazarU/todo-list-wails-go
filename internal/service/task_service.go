package service

import (
	"context"
	"errors"
	"wails-todo/internal/repo"
)

type TaskService struct {
	repo repo.TaskRepository
}

func NewTaskService(r repo.TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) CreateTask(ctx context.Context, t *repo.Task) error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	return s.repo.Create(ctx, t)
}

func (s *TaskService) GetAllTasks(ctx context.Context) ([]*repo.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*repo.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) UpdateTask(ctx context.Context, t *repo.Task) error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	return s.repo.Update(ctx, t)
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}


