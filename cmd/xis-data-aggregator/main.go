package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"runtime/debug"
	"sync"
	"time"
	grpcapi "xis-data-aggregator/internal/api/grpc"

	"github.com/gin-gonic/gin"
	_ "xis-data-aggregator/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"xis-data-aggregator/internal/api/rest"
	"xis-data-aggregator/internal/metrics"
	"xis-data-aggregator/internal/mocks"
	"xis-data-aggregator/internal/models"
	"xis-data-aggregator/internal/repository"

	"github.com/golang/glog"

	"xis-data-aggregator/config"
	"xis-data-aggregator/internal/service"
	"xis-data-aggregator/pkg/utils"
)

var (
	// BuildVersion is set on compile time.
	BuildVersion string
)

// @title           XIS Data Aggregator API
// @version         1.0
// @description     This is the API for the XIS Data Aggregator service.
// @host      localhost:8080 // Or your actual host and port
// @BasePath  /api/v1 // Base path for your API endpoints
func main() {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorln(utils.PanicErrStr(r, debug.Stack(), "main"))
		}

		glog.Flush()
	}()

	flag.Parse()
	_ = flag.Set("log_dir", "c:\\TEMP") // todo rm

	glog.Infof("Build version %s", BuildVersion)

	// Read and fill config
	cfg, err := config.GetXisDataAggregatorConfig()
	if err != nil {
		glog.Fatalf("init fail, config.GetXisDataAggregatorConfig() error: %v", err)
	}

	// Update config from flags if set
	cfg.UpdateConfigFromFlags()

	// Connect DB
	repo, err := repository.NewRedisRepository()
	defer func(repo *repository.RedisRepository) {
		err := repo.Close()
		if err != nil {
			glog.Errorf("repository.Close() error: %v", err)
		}
	}(repo)
	glog.Infoln("Connected to DB")

	var dataService = service.NewDataService(repo)

	// Init channels
	inputPacks := make(chan *models.Pack)
	metricsChan := make(chan bool)
	stopChan := make(chan struct{})
	glog.Infoln("Channels created")

	/*
	 * Note: It is a good practice to additionally intercept system interrupts
	 * for correct system shutdown:
	 *
	 * sysChan := make(chan os.Signal)
	 * signal.Notify(sgnChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)
	 */

	// Init, start, stop async running
	var wg sync.WaitGroup
	defer wg.Wait() // Wait for all valuable goroutines to finish

	// Create and start metrics collector
	metricsCollector := metrics.Collector{
		ProcessingResult: metrics.ProcessingResult{},
		InputChannel:     metricsChan,
	}
	go metricsCollector.Start(&wg, cfg)
	wg.Add(1)
	glog.Infoln("Metrics collector started")

	// Start data consuming and processing.
	for i := 0; i < cfg.WorkersCount; i++ {
		go service.ProcessData(&wg, dataService, inputPacks, metricsChan)
		wg.Add(1)
	}

	inputPacksGenerator := mocks.InputPacksGenerator{
		Interval:   time.Duration(cfg.InputIntervalMs) * time.Millisecond,
		PackLength: cfg.PackLength,
		OutputChan: inputPacks,
		StopChan:   stopChan}

	go inputPacksGenerator.Start(cfg)
	glog.Infoln("Pack generator started")

	// Start gRPC server in a goroutine
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort)) // /api/v2
		if err != nil {
			glog.Fatalf("failed to listen: %v", err)
		}

		s := grpc.NewServer()
		grpcapi.RegisterDataServiceServer(s, dataService)

		glog.Infof("gRPC Server started at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			glog.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	gin.SetMode(gin.DebugMode) // todo: gin.SetMode(gin.ReleaseMode)
	h := rest.NewDataHandler(dataService)
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v1.GET("data/:id", h.GetByID)
	v1.GET("data", h.ListByTimeRange)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start REST server
	glog.Infof("REST Server starting on port %d", cfg.RestPort)
	err = r.Run(fmt.Sprintf(":%d", cfg.RestPort))
	if err != nil {
		glog.Fatalf("gin.Run() error: %v", err)
	}

	/*
	 * Note: usually runs in a separate routine and is called by external signals
	 * or system interrupts to correctly terminate the service
	 */
	glog.Infoln("Sending stop signal to generator...")
	close(stopChan) // Send signals to all readers.

}
