package api

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"xis-data-aggregator/internal/models"
	"xis-data-aggregator/pb"
)

func TestDataToProto(t *testing.T) {
	id1 := uuid.New()

	// Определяем тестовые случаи
	tests := []struct {
		name    string
		input   *models.Data
		want    *pb.Data
		wantErr bool
		errMsg  string // Ожидаемое сообщение об ошибке
	}{
		{
			name: "Successful conversion",
			input: &models.Data{
				ID:        id1,
				Timestamp: 1678886400,
				Max:       100,
			},
			want: &pb.Data{
				Id:        id1.String(),
				Timestamp: 1678886400,
				Max:       100,
			},
			wantErr: false,
		},
		{
			name:    "Nil input data",
			input:   nil,
			want:    nil,
			wantErr: true,
			errMsg:  "data is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DataToProto(tt.input)

			if tt.wantErr {
				assert.Error(t, err, "Expected an error but got none")
				assert.Contains(t, err.Error(), tt.errMsg, "Error message mismatch")
				assert.Nil(t, got, "Expected nil result when error occurs")
			} else {
				assert.NoError(t, err, "Expected no error but got one: %v", err)
				assert.NotNil(t, got, "Expected non-nil result")
				assert.Equal(t, tt.want.Id, got.Id, "ID mismatch")
				assert.Equal(t, tt.want.Timestamp, got.Timestamp, "Timestamp mismatch")
				assert.Equal(t, tt.want.Max, got.Max, "Max mismatch")
			}
		})
	}
}

func TestProtoToData(t *testing.T) {
	// Generate UUIDs for test cases
	validUUID1 := uuid.New()
	validUUID2 := uuid.New()

	// Define test cases
	tests := []struct {
		name    string
		input   *pb.Data
		want    *models.Data
		wantErr bool
		errMsg  string // Expected substring in the error message
	}{
		{
			name: "Successful conversion with valid UUID",
			input: &pb.Data{
				Id:        validUUID1.String(),
				Timestamp: 1678886400, // Unix timestamp in seconds
				Max:       100,
			},
			want: &models.Data{
				ID:        validUUID1,
				Timestamp: 1678886400,
				Max:       100,
			},
			wantErr: false,
		},
		{
			name: "Successful conversion with another valid UUID and negative Max",
			input: &pb.Data{
				Id:        validUUID2.String(),
				Timestamp: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), // Unix timestamp in milliseconds
				Max:       -5,
			},
			want: &models.Data{
				ID:        validUUID2,
				Timestamp: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC).UnixMilli(),
				Max:       -5,
			},
			wantErr: false,
		},
		{
			name:    "Nil input pb.Data",
			input:   nil,
			want:    nil,
			wantErr: true,
			errMsg:  "pb.Data is nil",
		},
		{
			name: "Max field at int32 min value",
			input: &pb.Data{
				Id:        validUUID1.String(),
				Timestamp: 1678886400,
				Max:       -2147483648, // Minimum int32 value
			},
			want: &models.Data{
				ID:        validUUID1,
				Timestamp: 1678886400,
				Max:       -2147483648,
			},
			wantErr: false,
		},
		{
			name: "Max field at int32 max value",
			input: &pb.Data{
				Id:        validUUID2.String(),
				Timestamp: 1678886400,
				Max:       2147483647, // Maximum int32 value
			},
			want: &models.Data{
				ID:        validUUID2,
				Timestamp: 1678886400,
				Max:       2147483647,
			},
			wantErr: false,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProtoToData(tt.input)

			if tt.wantErr {
				// Assert that an error occurred
				assert.Error(t, err, "Expected an error for test case: %s", tt.name)
				// Assert that the error message contains the expected substring
				assert.Contains(t, err.Error(), tt.errMsg, "Error message mismatch for test case: %s", tt.name)
				// Assert that the returned data is nil
				assert.Nil(t, got, "Expected nil data when error occurs for test case: %s", tt.name)
			} else {
				// Assert that no error occurred
				assert.NoError(t, err, "Did not expect an error but got one: %v for test case: %s", err, tt.name)
				// Assert that the returned data is not nil
				assert.NotNil(t, got, "Expected non-nil data for test case: %s", tt.name)
				// Assert that the converted fields match the expected values
				assert.Equal(t, tt.want.ID, got.ID, "ID mismatch for test case: %s", tt.name)
				assert.Equal(t, tt.want.Timestamp, got.Timestamp, "Timestamp mismatch for test case: %s", tt.name)
				assert.Equal(t, tt.want.Max, got.Max, "Max mismatch for test case: %s", tt.name)
			}
		})
	}
}
