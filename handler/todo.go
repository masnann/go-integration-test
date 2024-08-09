package handler

import (
	"fmt"
	"go-integration-test/model"
	"go-integration-test/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(service service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: service}
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	todo := new(model.Todo)
	if err := c.Bind(todo); err != nil {
		fmt.Println("Binding Error:", err) 
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err := h.todoService.CreateTodo(todo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, todo)
}
