package metrics

import (
	"github.com/golang/glog"
	"sync"
	"xis-data-aggregator/config"
)

type ProcessingResult struct {
	SuccessfullyCount int
	FailedCount       int
}

type Collector struct {
	ProcessingResult ProcessingResult
	InputChannel     <-chan bool // Note:atomic is more often used for simple metrics
}

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
