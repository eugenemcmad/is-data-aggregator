// Package utils provides utility functions for error handling and data manipulation.
// This package contains helper functions for creating formatted errors, handling panics,
// and managing slice operations.
package utils

import "testing"

// TestMaxValue tests the GetMaxValue function with various test cases.
// This test uses table-driven tests to verify the function works correctly
// with different input scenarios including empty slices, positive numbers,
// negative numbers, mixed numbers, single elements, and duplicate values.
//
// Test cases include:
//   - Empty slice: should return error
//   - Positive numbers: should find maximum positive value
//   - Negative numbers: should find maximum negative value (closest to zero)
//   - Mixed numbers: should find maximum value regardless of sign
//   - Single element: should return that element
//   - Duplicate max value: should return the maximum value even if duplicated
func TestMaxValue(t *testing.T) {
	// Test Cases struct defines the structure for each test case
	type testCase struct {
		name    string // Case name for identification
		slice   []int  // Target slice to test
		want    int    // Expected maximum value
		wantErr bool   // Whether an error is expected
		errMsg  string // Expected error message if error is expected
	}

	// Tables of test cases covering various scenarios
	tests := []testCase{
		{
			name:    "Empty slice",
			slice:   []int{},
			want:    0, // 0 - default for error.
			wantErr: true,
			errMsg:  "slice is empty",
		},
		{
			name:    "Positive numbers",
			slice:   []int{1, 5, 2, 8, 3},
			want:    8,
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "Negative numbers",
			slice:   []int{-10, -5, -20, -1, -3},
			want:    -1,
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "Mixed numbers",
			slice:   []int{-5, 0, 10, -2, 7},
			want:    10,
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "Single element",
			slice:   []int{42},
			want:    42,
			wantErr: false,
			errMsg:  "",
		},
		{
			name:    "Duplicate max value",
			slice:   []int{10, 5, 10, 2, 7},
			want:    10,
			wantErr: false,
			errMsg:  "",
		},
	}

	// Check by test cases.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetMaxValue(tc.slice)

			// Check errors.
			if (err != nil) != tc.wantErr {
				t.Errorf("GetMaxDataValue() got error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check errors message.
			if tc.wantErr && err.Error() != tc.errMsg {
				t.Errorf("GetMaxDataValue() got error message = %q, want error message %q", err.Error(), tc.errMsg)
				return
			}

			// Check normal result.
			if !tc.wantErr && got != tc.want {
				t.Errorf("GetMaxDataValue() got = %v, want %v", got, tc.want)
			}
		})
	}
}
