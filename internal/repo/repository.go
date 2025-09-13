package repo

import (
	"context"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Completed bool       `json:"completed"`
	CreatedAt time.Time  `json:"created_at"`
	DueDate   *time.Time `json:"due_date,omitempty"`
	Priority  int16      `json:"priority"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetAll(ctx context.Context) ([]*Task, error)
	GetByID(ctx context.Context, id string) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}

type PostgresTaskRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTaskRepository(pool *pgxpool.Pool) *PostgresTaskRepository {
	return &PostgresTaskRepository{pool: pool}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *Task) error {
	query := `INSERT INTO tasks (title, completed, created_at, due_date, priority) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	row := r.pool.QueryRow(ctx, query, task.Title, task.Completed, task.CreatedAt, task.DueDate, task.Priority)
	return row.Scan(&task.ID, &task.CreatedAt)
}

func (r *PostgresTaskRepository) GetAll(ctx context.Context) ([]*Task, error) {
	query := `SELECT id, title, completed, created_at, due_date, priority FROM tasks ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []*Task
	for rows.Next() {
		t := &Task{}
		err := rows.Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.DueDate, &t.Priority)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	query := `SELECT id, title, completed, created_at, due_date, priority FROM tasks WHERE id = $1`
	t := &Task{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.Title, &t.Completed, &t.CreatedAt, &t.DueDate, &t.Priority)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *Task) error {
	query := `UPDATE tasks SET title = $1, completed = $2, due_date = $3, priority = $4 WHERE id = $5`
	_, err := r.pool.Exec(ctx, query, task.Title, task.Completed, task.DueDate, task.Priority, task.ID)
	return err
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}


