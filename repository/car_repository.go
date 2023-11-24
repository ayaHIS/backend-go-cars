package repository

import "example/backend-go-cars/domain"

type CarRepository interface {
	ListCars() ([]domain.Car, error)
	AddCar(newCar domain.Car) error
	RentCar(registration string) error
	ReturnCar(registration string) error
}
