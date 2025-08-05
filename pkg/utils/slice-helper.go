// Package utils provides utility functions for error handling and data manipulation.
// This package contains helper functions for creating formatted errors, handling panics,
// and managing slice operations.
package utils

import "fmt"

// GetMaxValue finds the maximum value in a slice of integers.
// This function iterates through the slice and returns the largest integer value.
// If the slice is empty, it returns an error.
//
// Parameters:
//   - Data: a slice of integers to search for the maximum value
//
// Returns:
//   - int: the maximum value found in the slice
//   - error: an error if the slice is empty, nil otherwise
//
// Example:
//
//	max, err := GetMaxValue([]int{1, 5, 2, 8, 3})
//	if err != nil {
//	    // handle error
//	}
//	// max will be 8
//
// Example with empty slice:
//
//	max, err := GetMaxValue([]int{})
//	if err != nil {
//	    // err will be "slice is empty"
//	}
func GetMaxValue(Data []int) (int, error) {
	if len(Data) == 0 {
		return 0, fmt.Errorf("slice is empty")
	}

	maxVal := Data[0]

	for _, value := range Data {
		if value > maxVal {
			maxVal = value
		}
	}

	return maxVal, nil
}
