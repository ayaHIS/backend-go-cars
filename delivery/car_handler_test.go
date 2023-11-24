package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example/backend-go-cars/domain"
	"example/backend-go-cars/repository"

	"github.com/stretchr/testify/assert"
)

func TestCarHandler_ListCarsHandler(t *testing.T) {
	repo := repository.NewInMemoryCarRepository()
	handler := NewCarHandler(repo)

	req, err := http.NewRequest("GET", "/cars", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.ListCarsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var cars []domain.Car
	err = json.Unmarshal(rr.Body.Bytes(), &cars)
	assert.NoError(t, err)
	assert.Empty(t, cars)
}

func TestCarHandler_AddCarHandler(t *testing.T) {
	repo := repository.NewInMemoryCarRepository()
	handler := NewCarHandler(repo)

	newCar := domain.Car{
		Model:        "TestCar",
		Registration: "ABC123",
		Mileage:      0,
		IsRented:     false,
	}
	body, err := json.Marshal(newCar)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.AddCarHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var addedCar domain.Car
	err = json.Unmarshal(rr.Body.Bytes(), &addedCar)
	assert.NoError(t, err)
	assert.Equal(t, newCar, addedCar)

	req, err = http.NewRequest("POST", "/cars", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.AddCarHandler(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
}
