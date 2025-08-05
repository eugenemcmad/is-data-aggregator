// Package api provides data mapping utilities for converting between internal models and protobuf representations.
package api

import (
	"fmt"
	"xis-data-aggregator/internal/models"
	"xis-data-aggregator/pb"

	"github.com/google/uuid"
)

// note: `proto` format is used in both REST and GRPC

// DataToProto converts a models.Data struct to its protobuf representation (pb.Data).
// Returns an error if the input data is nil.
func DataToProto(data *models.Data) (*pb.Data, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	pbData := pb.Data{
		Id:        data.ID.String(), // Convert UUID to string for protobuf
		Timestamp: data.Timestamp,
		Max:       int32(data.Max),
	}

	return &pbData, nil
}

// ProtoToData converts a protobuf pb.Data struct to the internal models.Data struct.
// Returns an error if the input pb.Data is nil or if the ID cannot be parsed as a UUID.
func ProtoToData(pbData *pb.Data) (*models.Data, error) {
	var err error
	if pbData == nil {
		return nil, fmt.Errorf("pb.Data is nil")
	}

	data := models.Data{
		Timestamp: pbData.Timestamp,
		Max:       int(pbData.Max),
	}

	// Parse the string ID from protobuf into a UUID
	data.ID, err = uuid.Parse(pbData.Id)

	return &data, err
}
