package service

import (
	"fmt"
	"sync"
	"xis-data-aggregator/internal/models"

	"github.com/golang/glog"
)

var closeOnce sync.Once

func ProcessData(wg *sync.WaitGroup, ds *DataService, inputChan <-chan *models.Pack, metricsChan chan<- bool) {
	defer wg.Done()

	var err error
	for pack := range inputChan {
		err = ProcessPack(pack, ds, metricsChan)
		if err != nil {
			glog.Errorln("process data error:", err)
		}

	}

	closeOnce.Do(func() {
		close(metricsChan)
	})

}

func ProcessPack(pack *models.Pack, ds *DataService, metricsChan chan<- bool) error {

	// Try map pack to data
	var data, err = models.MapPackToData(pack)
	switch {
	case err != nil:
		metricsChan <- false
		return err
	case data == nil: // extremely unlikely, reservation from nil pointer exception
		metricsChan <- false
		return fmt.Errorf("data is nil")
	}

	//  Try save to DB
	err = ds.Put(data)
	if err != nil {
		metricsChan <- false
		return err
	}

	metricsChan <- true
	return nil
}
