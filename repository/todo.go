package repository

import (
	"database/sql"
	"go-integration-test/model"
)

type TodoRepository interface {
	CreateTodo(todo *model.Todo) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) CreateTodo(todo *model.Todo) error {
	_, err := r.db.Exec("INSERT INTO todos (title) VALUES (?)", todo.Title)
	if err != nil {
		return err
	}
	return nil
}
