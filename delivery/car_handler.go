package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example/backend-go-cars/domain"
	"example/backend-go-cars/repository"

	"github.com/gorilla/mux"
)

type CarHandler struct {
	usecase repository.CarRepository
}

func NewCarHandler(usecase repository.CarRepository) *CarHandler {
	return &CarHandler{usecase: usecase}
}

func (h *CarHandler) ListCarsHandler(w http.ResponseWriter, r *http.Request) {
	cars, err := h.usecase.ListCars()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing cars: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (h *CarHandler) AddCarHandler(w http.ResponseWriter, r *http.Request) {
	var newCar domain.Car
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
		return
	}

	err = h.usecase.AddCar(newCar)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding car: %v", err), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCar)
}

func (h *CarHandler) RentCarHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]

	err := h.usecase.RentCar(registration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error renting car: %v", err), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CarHandler) ReturnCarHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	registration := params["registration"]

	err := h.usecase.ReturnCar(registration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error returning car: %v", err), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}
