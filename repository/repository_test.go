package repository

import (
	"testing"

	"example/backend-go-cars/domain"

	"github.com/stretchr/testify/assert"
)

func TestCarRepository_ListCars(t *testing.T) {
	repo := NewInMemoryCarRepository()

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

func TestCarRepository_AddCar(t *testing.T) {
	repo := NewInMemoryCarRepository()

	newCar := domain.Car{
		Model:        "TestCar",
		Registration: "ABC123",
		Mileage:      0,
		IsRented:     false,
	}
	err := repo.AddCar(newCar)
	assert.NoError(t, err)

	err = repo.AddCar(newCar)
	assert.Error(t, err)
}

func TestCarRepository_RentCar(t *testing.T) {
	repo := NewInMemoryCarRepository()

	newCar := domain.Car{
		Model:        "TestCar",
		Registration: "ABC123",
		Mileage:      0,
		IsRented:     false,
	}
	err := repo.AddCar(newCar)
	assert.NoError(t, err)

	err = repo.RentCar("ABC123")
	assert.NoError(t, err)

	err = repo.RentCar("ABC123")
	assert.Error(t, err)

	err = repo.RentCar("XYZ456")
	assert.Error(t, err)
}

func TestCarRepository_ReturnCar(t *testing.T) {
	repo := NewInMemoryCarRepository()

	rentedCar := domain.Car{
		Model:        "TestCar",
		Registration: "ABC123",
		Mileage:      100,
		IsRented:     true,
	}
	err := repo.AddCar(rentedCar)
	assert.NoError(t, err)

	err = repo.ReturnCar("ABC123")
	assert.NoError(t, err)

	err = repo.ReturnCar("XYZ456")
	assert.Error(t, err)

	err = repo.ReturnCar("DEF789")
	assert.Error(t, err)
}
