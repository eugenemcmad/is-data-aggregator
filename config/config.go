// Package config provides configuration structures and functions for the XIS Data Aggregator service.
package config

import "flag"

// workersCount is the default number of workers for reading, aggregating, and saving to the database (tuned for weak test DB).
// metricsBatchSize is the default number of metrics to batch before processing.
// restPort is the default port for the REST API server.
// grpcPort is the default port for the gRPC server.
// inputIntervalMs is the default interval (in milliseconds) for input simulation (tuned for weak test DB).
// packLength is the default length of a data pack.
const (
	workersCount     = 5 // for weak test db
	metricsBatchSize = 10
	restPort         = 8080
	grpcPort         = 50051
	inputIntervalMs  = 555 // for weak test db
	packLength       = 10
)

// XisDataAggregatorConfig holds all configuration parameters for the XIS Data Aggregator service.
type XisDataAggregatorConfig struct {
	// WorkersCount is the number of workers for reading, aggregating, and saving to the database.
	WorkersCount int

	// RestPort is the port for the REST API server.
	RestPort int
	// GrpcPort is the port for the gRPC server.
	GrpcPort int

	// MetricsBatchSize is the number of metrics to batch before processing.
	MetricsBatchSize int

	// Simulation parameters
	// InputIntervalMs is the interval (in milliseconds) for input simulation.
	InputIntervalMs int
	// PackLength is the length of a data pack.
	PackLength int
}

// GetXisDataAggregatorConfig initializes and returns a default XisDataAggregatorConfig.
// This function provides a default config instead of using external tools like Consul.
func GetXisDataAggregatorConfig() (*XisDataAggregatorConfig, error) {

	config := XisDataAggregatorConfig{
		WorkersCount:     workersCount,
		RestPort:         restPort,
		GrpcPort:         grpcPort,
		MetricsBatchSize: metricsBatchSize,
		InputIntervalMs:  inputIntervalMs,
		PackLength:       packLength,
	}

	return &config, nil
}

// UpdateConfigFromFlags updates the configuration fields from command-line flags if provided.
// Only non-zero flag values will override the existing config values.
func (cfg *XisDataAggregatorConfig) UpdateConfigFromFlags() {
	var workersCount, metricsBatchSize, restPort, grpcPort, inputIntervalMs, packLength int

	flag.IntVar(&workersCount, "workersCount", 0, "workers count")
	if workersCount > 0 {
		cfg.WorkersCount = workersCount
	}

	flag.IntVar(&metricsBatchSize, "b", 0, "metrics batch size")
	if metricsBatchSize > 0 {
		cfg.MetricsBatchSize = metricsBatchSize
	}

	flag.IntVar(&restPort, "r", 0, "rest port")
	if restPort > 0 {
		cfg.RestPort = restPort
	}

	flag.IntVar(&grpcPort, "g", 0, "grpc port")
	if grpcPort > 0 {
		cfg.GrpcPort = grpcPort
	}

	flag.IntVar(&inputIntervalMs, "n", 0, "input interval")
	if inputIntervalMs > 0 {
		cfg.InputIntervalMs = inputIntervalMs
	}

	flag.IntVar(&packLength, "l", 0, "input pack length")
	if packLength > 0 {
		cfg.PackLength = packLength
	}
}
