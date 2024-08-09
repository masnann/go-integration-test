package service

import (
	"go-integration-test/model"
	"go-integration-test/repository"
)

type TodoService struct {
	todoRepo repository.TodoRepository
}

func NewTodoService(todoRepo repository.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (s *TodoService) CreateTodo(todo *model.Todo) error {
	err := s.todoRepo.CreateTodo(todo)
	if err != nil {
		return err
	}

	return nil
}
