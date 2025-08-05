package mocks

import (
	"errors"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"math/rand"
	"time"
	"xis-data-aggregator/config"
	"xis-data-aggregator/internal/models"
)

const valueLimit = 1000 // limitations for ease of interpretation of results

type InputPacksGenerator struct {
	Interval   time.Duration
	PackLength int
	OutputChan chan<- *models.Pack // Channel for output
	StopChan   chan struct{}       // Channel for stopping
}

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

func GeneratePack(dataLength int) (*models.Pack, error) {
	if dataLength <= 0 {
		return nil, errors.New("dataLength must be a positive integer")
	}

	// Generate uniq UUID
	id := uuid.New()

	// Set Timestamp (Unix nanoseconds for high precision)
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
