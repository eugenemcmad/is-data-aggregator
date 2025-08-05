// Package mocks provides mock data generators for testing and development purposes.
package mocks

import (
	"errors"
	"math/rand"
	"time"
	"xis-data-aggregator/config"
	"xis-data-aggregator/internal/models"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

// valueLimit defines the upper bound for generated random values in packs.
const valueLimit = 1000 // limitations for ease of interpretation of results

// InputPacksGenerator generates mock Pack data at regular intervals and sends them to an output channel.
type InputPacksGenerator struct {
	Interval   time.Duration       // Time interval between generated packs
	PackLength int                 // Number of data points in each generated pack
	OutputChan chan<- *models.Pack // Channel for outputting generated packs
	StopChan   chan struct{}       // Channel to signal generator to stop
}

// Start begins generating packs at the specified interval until StopChan is closed.
// Packs are sent to OutputChan. Logs errors and debug info as appropriate.
func (g *InputPacksGenerator) Start(cfg *config.XisDataAggregatorConfig) {
	ticker := time.NewTicker(g.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pack, err := GeneratePack(g.PackLength)
			if err != nil {
				glog.Errorln("Error creating pack: %s", err)
				continue
			}

			g.OutputChan <- pack

			if rand.Int()%cfg.MetricsBatchSize == 0 {
				glog.Infof("dbg: generated pack: %v\n", *pack)
			}

		case <-g.StopChan:
			glog.Infoln("Pack generator stopping.")
		}
	}
}

// GeneratePack creates a new Pack with a unique ID, current timestamp, and a slice of random integers.
// Returns an error if dataLength is not positive.
func GeneratePack(dataLength int) (*models.Pack, error) {
	if dataLength <= 0 {
		return nil, errors.New("dataLength must be a positive integer")
	}

	// Generate uniq UUID for the pack
	id := uuid.New()

	// Set Timestamp (Unix microseconds for high precision)
	timestamp := time.Now().UnixMicro()

	// Generate Data slice of random integers
	data := make([]int, dataLength)
	for i := 0; i < dataLength; i++ {
		data[i] = rand.Intn(valueLimit)
	}

	// Create and return the Pack instance
	return &models.Pack{
		ID:        id,
		Timestamp: timestamp,
		Data:      data,
	}, nil
}
