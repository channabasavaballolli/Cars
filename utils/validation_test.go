package utils

import (
	"car-service/models"
	"testing"
)

func TestValidateCar(t *testing.T) {
	// Test Valid Car
	validCar := models.Car{Price: 100, Year: 2020}
	if err := ValidateCar(validCar); err != nil {
		t.Errorf("Expected valid car, got error: %v", err)
	}

	// Test Invalid Price
	invalidPriceCar := models.Car{Price: 0, Year: 2020}
	if err := ValidateCar(invalidPriceCar); err == nil {
		t.Error("Expected error for price 0, got nil")
	}

	// Test Invalid Year
	invalidYearCar := models.Car{Price: 100, Year: 1800}
	if err := ValidateCar(invalidYearCar); err == nil {
		t.Error("Expected error for year 1800, got nil")
	}
}
