// Package grpc provides the gRPC server implementation for the DataService.
package grpc

import (
	"errors"
	"io"
	"strconv"
	"xis-data-aggregator/internal/repository"

	"xis-data-aggregator/internal/api"
	"xis-data-aggregator/internal/service"
	"xis-data-aggregator/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

// DataServiceServer implements the gRPC DataService interface.
type DataServiceServer struct {
	pb.UnimplementedDataServiceServer                      // Embeds unimplemented server for forward compatibility
	service                           *service.DataService // Business logic service
}

// NewDataServiceServer creates a new gRPC DataServiceServer instance.
// Takes a pointer to the business logic service.
func NewDataServiceServer(service *service.DataService) *DataServiceServer {
	return &DataServiceServer{
		service: service,
	}
}

// RegisterDataServiceServer registers the DataServiceServer with the given gRPC server.
func RegisterDataServiceServer(s *grpc.Server, service *service.DataService) {
	server := NewDataServiceServer(service)
	pb.RegisterDataServiceServer(s, server)
}

// GetDataById handles bidirectional streaming for getting data by ID.
// Receives requests with IDs from the client, fetches data, and streams responses back.
// Returns a gRPC error if the request is invalid or if data is not found.
func (s *DataServiceServer) GetDataById(stream pb.DataService_GetDataByIdServer) error {
	glog.Infoln("GetDataById stream started")
	defer glog.Infoln("GetDataById stream ended")

	for {
		// Receive request from client
		req, err := stream.Recv()
		if err == io.EOF {
			glog.Infoln("Client closed stream")
			return nil
		}
		if err != nil {
			glog.Errorf("Error receiving request: %v", err)
			return status.Errorf(codes.Internal, "failed to receive request: %v", err)
		}

		if req == nil || req.Id == "" {
			glog.Errorln("Invalid request: ID is empty")
			return status.Errorf(codes.InvalidArgument, "ID cannot be empty")
		}

		// Parse UUID from request
		id, err := uuid.Parse(req.Id)
		if err != nil {
			glog.Errorf("Invalid UUID format: %v", err)
			return status.Errorf(codes.InvalidArgument, "invalid UUID format: %v", err)
		}

		// Get data from service layer
		data, err := s.service.GetByID(id)

		switch {
		case errors.Is(err, repository.ErrNotFound):
			glog.Infof("Data not found for ID: %s", req.Id)
			return status.Errorf(codes.NotFound, "data not found for ID: %s", req.Id)
		case errors.Is(err, service.ErrNotFound):
			glog.Infof("Data not found for ID: %s", req.Id)
			return status.Errorf(codes.NotFound, "data not found for ID: %s", req.Id)
		case err != nil:
			glog.Errorf("Service error: %v", err)
			return status.Errorf(codes.Internal, "internal server error: %v", err)
		}

		// Convert to proto format for response
		protoData, err := api.DataToProto(data)
		if err != nil {
			glog.Errorf("Error converting data to proto: %v", err)
			return status.Errorf(codes.Internal, "failed to convert data: %v", err)
		}

		// Send response back to client
		if err := stream.Send(protoData); err != nil {
			glog.Errorf("Error sending response: %v", err)
			return status.Errorf(codes.Internal, "failed to send response: %v", err)
		}

		glog.Infof("Successfully sent data for ID: %s", req.Id)
	}
}

// ListDataByTimeRange handles bidirectional streaming for listing data by time range.
// Receives requests with time range parameters, fetches matching data, and streams responses back.
// Returns a gRPC error if the request is invalid or if no data is found.
func (s *DataServiceServer) ListDataByTimeRange(stream pb.DataService_ListDataByTimeRangeServer) error {
	glog.Infoln("ListDataByTimeRange stream started")
	defer glog.Infoln("ListDataByTimeRange stream ended")

	for {
		// Receive request from client
		req, err := stream.Recv()
		if err == io.EOF {
			glog.Infoln("Client closed stream")
			return nil
		}
		if err != nil {
			glog.Errorf("Error receiving request: %v", err)
			return status.Errorf(codes.Internal, "failed to receive request: %v", err)
		}

		if req == nil {
			glog.Errorln("Request is nil")
			return status.Errorf(codes.InvalidArgument, "request cannot be nil")
		}

		// Parse time range parameters from request
		from, err := strconv.ParseInt(req.From, 10, 64)
		if err != nil {
			glog.Errorf("Invalid 'from' parameter: %v", err)
			return status.Errorf(codes.InvalidArgument, "invalid 'from' parameter: %v", err)
		}

		to, err := strconv.ParseInt(req.To, 10, 64)
		if err != nil {
			glog.Errorf("Invalid 'to' parameter: %v", err)
			return status.Errorf(codes.InvalidArgument, "invalid 'to' parameter: %v", err)
		}

		// Validate time range
		if from >= to {
			glog.Errorln("Invalid time range: 'from' must be less than 'to'")
			return status.Errorf(codes.InvalidArgument, "invalid time range: 'from' must be less than 'to'")
		}

		// Get data from service layer for the specified period
		dataList, err := s.service.ListByPeriod(from, to)

		switch {
		case errors.Is(err, repository.ErrNotFound):
			glog.Infof("No data found for time range: %d to %d", from, to)
			return status.Errorf(codes.NotFound, "no data found for time range: %d to %d", from, to)
		case errors.Is(err, service.ErrNotFound):
			glog.Infof("No data found for time range: %d to %d", from, to)
			return status.Errorf(codes.NotFound, "no data found for time range: %d to %d", from, to)
		case err != nil:
			glog.Errorf("Service error: %v", err)
			return status.Errorf(codes.Internal, "internal server error: %v", err)
		}

		// Convert data list to proto format for response
		protoDataList := make([]*pb.Data, len(dataList))
		for i, data := range dataList {
			protoData, err := api.DataToProto(&data)
			if err != nil {
				glog.Errorf("Error converting data to proto: %v", err)
				return status.Errorf(codes.Internal, "failed to convert data: %v", err)
			}
			protoDataList[i] = protoData
		}

		// Create response message
		response := &pb.ListDataByTimeRangeResponse{
			DataItems: protoDataList,
		}

		// Send response back to client
		if err := stream.Send(response); err != nil {
			glog.Errorf("Error sending response: %v", err)
			return status.Errorf(codes.Internal, "failed to send response: %v", err)
		}

		glog.Infof("Successfully sent %d data items for time range: %d to %d", len(dataList), from, to)
	}
}
