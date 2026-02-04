package utils

import "testing"

func TestIsPriceValid(t *testing.T) {
	if IsPriceValid(100) != true {
		t.Error("Expected 100 to be valid, but it was false")
	}
	if IsPriceValid(-50) != false {
		t.Error("Expected -50 to be invalid, but it was true")
	}
	if IsPriceValid(0) != false {
		t.Error("Expected 0 to be invalid, but it was true")
	}
}
