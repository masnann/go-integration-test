package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-integration-test/model"
	"go-integration-test/repository/mocks"
	"go-integration-test/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTodo(t *testing.T) {
	// Persiapkan mock repository
	mocksRepo := mocks.NewTodoRepository(t)
	expectedTodo := &model.Todo{ID: 1, Title: "Test Todo"}
	mocksRepo.On("CreateTodo", mock.AnythingOfType("*model.Todo")).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*model.Todo)
		arg.ID = expectedTodo.ID
		arg.Title = expectedTodo.Title
	}).Return(nil).Once()

	// Buat instance service dan handler
	todoService := service.NewTodoService(mocksRepo)
	handler := NewTodoHandler(*todoService)

	// Buat request dengan payload
	payload := map[string]string{"title": "Test Todo"}
	payloadBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(payloadBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	if err := handler.CreateTodo(c); err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Periksa status code
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Periksa response body
	var responseTodo model.Todo
	if err := json.NewDecoder(rec.Body).Decode(&responseTodo); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedTodo.ID, responseTodo.ID)
	assert.Equal(t, expectedTodo.Title, responseTodo.Title)

	mocksRepo.AssertExpectations(t)
}

func TestCreateTodo_BindingError(t *testing.T) {
	// Persiapkan mock repository
	mocksRepo := mocks.NewTodoRepository(t)

	// Buat request dengan payload yang tidak valid
	payload := `{"title": 1}` // JSON tidak valid
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader([]byte(payload)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	// Buat instance service dan handler
	todoService := service.NewTodoService(mocksRepo)
	handler := NewTodoHandler(*todoService)

	// Jalankan handler
	if err := handler.CreateTodo(c); err == nil {
		t.Errorf("Expected error, got nil")
		return
	}

	// // Periksa status code
	// if rec.Code != http.StatusBadRequest {
	//     t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rec.Code)
	// }

	// Periksa body response untuk debugging
	body := rec.Body.String()
	t.Logf("Response Body: %s", body)
}

func TestCreateTodo_ServiceError(t *testing.T) {
    // Persiapkan mock repository
    mocksRepo := mocks.NewTodoRepository(t)
    expectedError := errors.New("service error")

    // Atur mock repository untuk mengembalikan error saat CreateTodo dipanggil
    mocksRepo.On("CreateTodo", mock.AnythingOfType("*model.Todo")).Return(expectedError).Once()

    // Buat request dengan payload yang valid
    payload := map[string]string{"title": "Test Todo"}
    payloadBytes, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(payloadBytes))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    e := echo.New()
    c := e.NewContext(req, rec)

    // Buat instance service dan handler
    todoService := service.NewTodoService(mocksRepo)
    handler := NewTodoHandler(*todoService)

    // Jalankan handler
    err := handler.CreateTodo(c)
    if err == nil {
        t.Errorf("Expected error, got nil")
        return
    }

    // // Periksa status code
    // if rec.Code != http.StatusInternalServerError {
    //     t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
    // }

    // Periksa body response untuk debugging
    body := rec.Body.String()
    t.Logf("Response Body: %s", body)

    // Periksa apakah error yang diharapkan dikembalikan
    var httpError *echo.HTTPError
    if !errors.As(err, &httpError) {
        t.Fatalf("Expected echo.HTTPError, got %v", err)
    }
    assert.Equal(t, http.StatusInternalServerError, httpError.Code)
    assert.Contains(t, httpError.Message, expectedError.Error())
}
