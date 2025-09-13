package domain

import (
	"context"
	"wails-todo/internal/repo"
	"wails-todo/internal/service"
)

type TaskUsecase struct {
	service *service.TaskService
}

func NewTaskUsecase(s *service.TaskService) *TaskUsecase {
	return &TaskUsecase{service: s}
}

func (u *TaskUsecase) CreateTask(ctx context.Context, t *repo.Task) error {
	return u.service.CreateTask(ctx, t)
}

func (u *TaskUsecase) GetAllTasks(ctx context.Context) ([]*repo.Task, error) {
	return u.service.GetAllTasks(ctx)
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, t *repo.Task) error {
	return u.service.UpdateTask(ctx, t)
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id string) error {
	return u.service.DeleteTask(ctx, id)
}


