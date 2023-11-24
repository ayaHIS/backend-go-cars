package repository

import (
	"fmt"
	"sync"
	"testing"

	"example/backend-go-cars/domain"

	"github.com/stretchr/testify/assert"
)

type InMemoryTestCarRepository struct {
	cars []domain.Car
	mu   sync.Mutex
}

func NewInMemoryTestCarRepository() *InMemoryTestCarRepository {
	return &InMemoryTestCarRepository{}
}
func (r *InMemoryTestCarRepository) ListCars() ([]domain.Car, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.cars, nil
}

func (r *InMemoryTestCarRepository) AddCar(newCar domain.Car) error {
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

func (r *InMemoryTestCarRepository) RentCar(registration string) error {
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

func (r *InMemoryTestCarRepository) ReturnCar(registration string) error {
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

func (r *InMemoryTestCarRepository) ResetState() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cars = nil
}

func TestInMemoryTestCarRepository_ListCars(t *testing.T) {
	repo := NewInMemoryTestCarRepository()

	cars, err := repo.ListCars()
	assert.NoError(t, err)
	assert.Empty(t, cars)

	newCar := domain.Car{
		Model:        "TestCar",
		Registration: "ABC123",
		Mileage:      0,
		IsRented:     false,
	}
	err = repo.AddCar(newCar)
	assert.NoError(t, err)
	cars, err = repo.ListCars()
	assert.NoError(t, err)
	assert.Len(t, cars, 1)
	assert.Equal(t, newCar, cars[0])
}
