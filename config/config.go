package config

import "flag"

const (
	workersCount     = 5 // for weak test db
	metricsBatchSize = 10
	restPort         = 8080
	grpcPort         = 50051
	inputIntervalMs  = 555 // for weak test db
	packLength       = 10
)

type XisDataAggregatorConfig struct {
	//SqliteConfig
	//RedisConfig

	WorkersCount int // Count of workers for read, aggregate and save to db.

	RestPort int
	GrpcPort int

	MetricsBatchSize int

	// Simulation
	InputIntervalMs int
	PackLength      int
}

// GetXisDataAggregatorConfig default config init instead of such tools as Consul for example
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
