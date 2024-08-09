package main

import (
	"database/sql"
	"go-integration-test/handler"
	"go-integration-test/repository"
	"go-integration-test/service"

	"github.com/labstack/echo/v4"
)

func main() {
	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(*todoService)

	e := echo.New()
	e.POST("/todos", todoHandler.CreateTodo)

	e.Logger.Fatal(e.Start(":1323"))
}
