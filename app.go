package main

import (
	"context"
	"time"
	"wails-todo/internal/domain"
	"wails-todo/internal/repo"
	"wails-todo/internal/service"
)

// App представляет основную структуру Wails приложения
// Содержит контекст и бизнес-логику для работы с задачами
type App struct {
	ctx     context.Context
	usecase *domain.TaskUsecase
}

// NewApp создает новый экземпляр приложения
func NewApp() *App {
	return &App{}
}

// startup инициализирует приложение при запуске
// Настраивает зависимости и подключает репозиторий
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// Используем in-memory репозиторий для демонстрации
	// В продакшене здесь был бы PostgreSQL
	repository := repo.NewMemoryTaskRepository()
	service := service.NewTaskService(repository)
	a.usecase = domain.NewTaskUsecase(service)
}

// ListTasks возвращает все задачи из репозитория
// Вызывается фронтендом для отображения списка задач
func (a *App) ListTasks() ([]*repo.Task, error) {
	return a.usecase.GetAllTasks(a.ctx)
}

// CreateTask создает новую задачу
// Парсит дату и приоритет из строковых параметров фронтенда
func (a *App) CreateTask(title string, description string, dueISO string, priority string) error {
	var due *time.Time
	if dueISO != "" {
		t, err := time.Parse(time.RFC3339, dueISO)
		if err != nil {
			return err
		}
		due = &t
	}
	
	// Конвертируем строковый приоритет в числовой
	prio := int16(2) // medium по умолчанию
	if priority == "low" {
		prio = 1
	} else if priority == "high" {
		prio = 3
	}
	
	task := &repo.Task{
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
		DueDate:   due,
		Priority:  prio,
	}
	return a.usecase.CreateTask(a.ctx, task)
}

// DeleteTask удаляет задачу по ID
func (a *App) DeleteTask(id string) error {
	return a.usecase.DeleteTask(a.ctx, id)
}

// UpdateTask обновляет существующую задачу
// Используется для изменения статуса выполнения
func (a *App) UpdateTask(task *repo.Task) error {
	return a.usecase.UpdateTask(a.ctx, task)
}
