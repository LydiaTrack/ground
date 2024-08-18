package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Generate6DigitCode generates a random 6-digit number as a string.
// If includeZero is false, the number will be between 100000 and 999999 (inclusive).
func Generate6DigitCode(includeZero bool) (string, error) {
	var n *big.Int
	var err error

	if includeZero {
		// Define the maximum value for a 6-digit number (exclusive)
		maxInt := big.NewInt(1000000) // 1 million
		// Generate a random number between 0 and 999999
		n, err = rand.Int(rand.Reader, maxInt)
		if err != nil {
			return "", fmt.Errorf("failed to generate a random number: %v", err)
		}
	} else {
		// Define the minimum and maximum values for a 6-digit number (inclusive)
		minInt := big.NewInt(100000) // 100000 (inclusive)
		maxInt := big.NewInt(900000) // 900000 (range)
		// Generate a random number between 100000 and 999999
		randRange, err := rand.Int(rand.Reader, maxInt)
		if err != nil {
			return "", fmt.Errorf("failed to generate a random number: %v", err)
		}
		n = randRange.Add(randRange, minInt)
	}

	// Format the number as a 6-digit string
	return fmt.Sprintf("%06d", n.Int64()), nil
}
