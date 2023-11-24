package repository

import (
	"fmt"
	"sync"

	"example/backend-go-cars/domain"
)

type InMemoryCarRepository struct {
	cars []domain.Car
	mu   sync.Mutex
}

func NewInMemoryCarRepository() *InMemoryCarRepository {
	return &InMemoryCarRepository{}
}

func (r *InMemoryCarRepository) ListCars() ([]domain.Car, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cars, nil
}

func (r *InMemoryCarRepository) AddCar(newCar domain.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, car := range r.cars {
		if car.Registration == newCar.Registration {
			return fmt.Errorf("car with registration '%s' already exists", newCar.Registration)
		}
	}

	r.cars = append(r.cars, newCar)
	return nil
}

func (r *InMemoryCarRepository) RentCar(registration string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, car := range r.cars {
		if car.Registration == registration {
			if car.IsRented {
				return fmt.Errorf("car with registration '%s' is already rented", registration)
			}
			r.cars[i].IsRented = true
			return nil
		}
	}

	return fmt.Errorf("car with registration '%s' not found", registration)
}

func (r *InMemoryCarRepository) ReturnCar(registration string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, car := range r.cars {
		if car.Registration == registration {
			if !car.IsRented {
				return fmt.Errorf("car with registration '%s' is not marked as rented", registration)
			}

			r.cars[i].Mileage += 100
			r.cars[i].IsRented = false
			return nil
		}
	}

	return fmt.Errorf("car with registration '%s' not found", registration)
}
