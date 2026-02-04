package utils

import "errors"

func IsPriceValid(price float64) bool {
	if price <= 0 {
		return false
	}
	return true
}

func ValidateCar(price float64, year int) error {
	if !IsPriceValid(price) {
		return errors.New("price must be greater than zero")
	}
	if year < 1886 { // First car in 1886
		return errors.New("year cannot be before 1886")
	}
	return nil
}
