package decode

import (
	"math"
)

// CToF converts a temperature from Celsius to Fahrenheit, rounded to one
// decimal digit.
func CToF(tempC float64) float64 {
	tempF := tempC*9/5 + 32

	return math.Round(tempF*10) / 10
}
