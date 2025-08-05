// Package metrics provides simple in-memory metrics collection for processing results.
package metrics

import (
	"sync"
	"xis-data-aggregator/config"

	"github.com/golang/glog"
)

// ProcessingResult holds counters for successful and failed processing attempts.
type ProcessingResult struct {
	// SuccessfullyCount is the number of successfully processed items.
	SuccessfullyCount int
	// FailedCount is the number of failed processing attempts.
	FailedCount int
}

// Collector collects processing metrics from an input channel.
type Collector struct {
	// ProcessingResult stores the current counts of successes and failures.
	ProcessingResult ProcessingResult
	// InputChannel receives boolean values indicating processing success (true) or failure (false).
	InputChannel <-chan bool // Note:atomic is more often used for simple metrics
}

// Start begins collecting metrics from the InputChannel and logs batch results.
// It should be run as a goroutine and will signal completion on the provided WaitGroup.
// Metrics are logged every cfg.MetricsBatchSize successful or failed items.
func (o *Collector) Start(wg *sync.WaitGroup, cfg *config.XisDataAggregatorConfig) {
	defer wg.Done()
	// Use log for metrics instead of tools as prometheus for example
	for ok := range o.InputChannel {
		if ok {
			o.ProcessingResult.SuccessfullyCount++
			if o.ProcessingResult.SuccessfullyCount%cfg.MetricsBatchSize == 0 {
				glog.Infof("Successfully processed %d items\n", o.ProcessingResult.SuccessfullyCount)
			}
		} else {
			o.ProcessingResult.FailedCount++
			if o.ProcessingResult.SuccessfullyCount%cfg.MetricsBatchSize == 0 {
				glog.Infof("Failed processing %d items\n", o.ProcessingResult.FailedCount)
			}
		}
	}
}
