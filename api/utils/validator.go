package utils

import (
	validator "github.com/asaskevich/govalidator"
)

// IsInteger is a func to validate if all value is integer from string type values
func IsInteger(values ...string) bool {
	for _, value := range values {
		if !validator.IsInt(value) {
			return false
		}
	}

	return true
}

// IsNonNegative is a func to validate if all value is non-negative integer from integer type values
func IsNonNegative(values ...int) bool {
	for _, value := range values {
		if value < 0 {
			return false
		}
	}

	return true
}
