package utils

import "testing"

func TestMaxValue(t *testing.T) {
	// Test Cases struct
	type testCase struct {
		name    string // Case name
		slice   []int  // Target object
		want    int    // Expected value
		wantErr bool   // Error expected
		errMsg  string // Expected error message
	}

	// Tables of test cases
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
