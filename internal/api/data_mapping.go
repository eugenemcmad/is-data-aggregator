package api

import (
	"fmt"
	"github.com/google/uuid"
	"xis-data-aggregator/internal/models"
	"xis-data-aggregator/pb"
)

// note: `proto` format is used in both REST and GRPC

func DataToProto(data *models.Data) (*pb.Data, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	pbData := pb.Data{
		Id:        data.ID.String(),
		Timestamp: data.Timestamp,
		Max:       int32(data.Max),
	}

	return &pbData, nil
}

func ProtoToData(pbData *pb.Data) (*models.Data, error) {
	var err error
	if pbData == nil {
		return nil, fmt.Errorf("pb.Data is nil")
	}

	data := models.Data{
		Timestamp: pbData.Timestamp,
		Max:       int(pbData.Max),
	}

	data.ID, err = uuid.Parse(pbData.Id)

	return &data, err
}
