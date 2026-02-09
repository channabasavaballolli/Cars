package utils

import (
	"car-service/models"
	"errors"
)

func ValidateCar(car models.Car) error {
	if car.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if car.Year <= 1886 {
		return errors.New("year must be greater than 1886")
	}
	return nil
}
