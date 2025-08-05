# gRPC Handlers Implementation

This document describes how to create, initialize, and start gRPC handlers for the XIS Data Aggregator service.

## Overview

The gRPC service provides bidirectional streaming endpoints for data retrieval operations. The service is implemented using the `google.golang.org/grpc` package and follows the standard gRPC patterns.

## Service Definition

The service is defined in `gen/proto/data.proto` and includes two main methods:

1. **GetDataById** - Retrieves data by UUID
2. **ListDataByTimeRange** - Retrieves data within a specified time range

Both methods use bidirectional streaming for request/response handling.

## Implementation Structure

### 1. Server Implementation (`internal/api/grpc/data_server.go`)

The gRPC server implementation includes:

- **DataServiceServer** - Main server struct that implements the gRPC interface
- **NewDataServiceServer** - Constructor function
- **RegisterDataServiceServer** - Registration function for the gRPC server
- **GetDataById** - Handler for retrieving data by ID
- **ListDataByTimeRange** - Handler for retrieving data by time range

### 2. Key Features

- **Error Handling**: Proper gRPC status codes and error messages
- **Logging**: Comprehensive logging using glog
- **Data Conversion**: Automatic conversion between internal models and protobuf messages
- **Validation**: Input validation for UUIDs and time ranges
- **Streaming**: Bidirectional streaming support

## Usage

### Starting the Server

The gRPC server is automatically started alongside the REST server in `main.go`:

```go
// Start gRPC server in a goroutine
go func() {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
    if err != nil {
        glog.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    grpcapi.RegisterDataServiceServer(s, svc)

    glog.Infof("gRPC Server started at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        glog.Fatalf("failed to serve gRPC: %v", err)
    }
}()
```

### Client Example

See `examples/grpc_client_example.go` for a complete client implementation.

## API Endpoints

### GetDataById

**Request:**
```protobuf
message GetDataByIDRequest {
    string id = 1;
}
```

**Response:**
```protobuf
message Data {
    string id = 1;
    int64 timestamp = 2;
    int32 max = 3;
}
```

**Usage:**
1. Create a bidirectional stream
2. Send a request with a valid UUID
3. Receive the corresponding data
4. Close the stream

### ListDataByTimeRange

**Request:**
```protobuf
message ListDataByTimeRangeRequest {
    string from = 1;
    string to = 2;
}
```

**Response:**
```protobuf
message ListDataByTimeRangeResponse {
    repeated Data data_items = 1;
}
```

**Usage:**
1. Create a bidirectional stream
2. Send a request with time range parameters (Unix timestamps as strings)
3. Receive a list of data items
4. Close the stream

## Error Handling

The server returns appropriate gRPC status codes:

- **InvalidArgument**: Invalid UUID format or time range
- **NotFound**: Data not found for the given criteria
- **Internal**: Server errors or data conversion issues

## Configuration

The gRPC server port is configured via the `GrpcPort` field in the configuration. Default is typically 50051.

## Testing

To test the gRPC server:

1. Start the server: `go run cmd/xis-data-aggregator/main.go`
2. Run the client example: `go run examples/grpc_client_example.go`

## Dependencies

Required dependencies in `go.mod`:
- `google.golang.org/grpc`
- `google.golang.org/protobuf`
- `github.com/google/uuid`
- `github.com/golang/glog`

## Best Practices

1. **Error Handling**: Always check for errors and return appropriate gRPC status codes
2. **Logging**: Use structured logging for debugging and monitoring
3. **Resource Management**: Properly close streams and connections
4. **Validation**: Validate all input parameters before processing
5. **Performance**: Use streaming for large datasets
6. **Security**: Implement authentication and authorization as needed 