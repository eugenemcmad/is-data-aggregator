package main

import (
	"flag"
	"fmt"
	"net"
	"runtime/debug"
	"sync"
	"time"
	"xis-data-aggregator/config"
	_ "xis-data-aggregator/docs"
	grpcapi "xis-data-aggregator/internal/api/grpc"
	"xis-data-aggregator/internal/api/rest"
	"xis-data-aggregator/internal/metrics"
	"xis-data-aggregator/internal/mocks"
	"xis-data-aggregator/internal/models"
	"xis-data-aggregator/internal/repository"
	"xis-data-aggregator/internal/service"
	"xis-data-aggregator/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

// @title           XIS Data Aggregator API
// @version         1.0
// @description     This is the API for the XIS Data Aggregator service.
// @host      localhost:8080 // Or your actual host and port
// @BasePath  /api/v1 // Base path for your API endpoints
// Entry point for the XIS Data Aggregator service
func main() {
	// Defer a panic handler to log any unexpected errors and flush logs on exit
	defer func() {
		if r := recover(); r != nil {
			// Note: Use `glog` without parameters. Default directory is ./tmp.
			glog.Errorln(utils.PanicErrStr(r, debug.Stack(), "main"))
		}

		glog.Flush() // Flush logs.
	}()

	flag.Parse() // Parse command-line flags

	// Read configuration from file/environment and handle errors
	cfg, err := config.GetXisDataAggregatorConfig()
	if err != nil {
		glog.Fatalf("init fail, config.GetXisDataAggregatorConfig() error: %v", err)
	}

	// Optionally override config values with command-line flags
	cfg.UpdateConfigFromFlags()

	// Initialize Redis repository (database connection)
	repo, err := repository.NewRedisRepository()
	defer func(repo *repository.RedisRepository) {
		err := repo.Close()
		if err != nil {
			glog.Errorf("repository.Close() error: %v", err)
		}
	}(repo)
	glog.Infoln("Connected to DB")

	// Create the main data service with the repository
	var dataService = service.NewDataService(repo)

	// Initialize channels for inter-goroutine communication
	inputPacks := make(chan *models.Pack)
	metricsChan := make(chan bool)
	stopChan := make(chan struct{})
	glog.Infoln("Channels created")

	/*
	 * Note: It is a good practice to additionally intercept system interrupts
	 * for correct system shutdown:
	 *
	 * sigChan := make(chan os.Signal)
	 * signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	 */

	// Use a WaitGroup to manage goroutines and ensure clean shutdown
	var wg sync.WaitGroup
	defer wg.Wait() // Wait for all valuable goroutines to finish

	// Create and start the metrics collector goroutine
	metricsCollector := metrics.Collector{
		ProcessingResult: metrics.ProcessingResult{},
		InputChannel:     metricsChan,
	}
	go metricsCollector.Start(&wg, cfg)
	wg.Add(1)
	glog.Infoln("Metrics collector started")

	// Start worker goroutines for data processing
	for i := 0; i < cfg.WorkersCount; i++ {
		go service.ProcessData(&wg, dataService, inputPacks, metricsChan)
		wg.Add(1)
	}

	// Create and start the mock input pack generator (simulates incoming data)
	inputPacksGenerator := mocks.InputPacksGenerator{
		Interval:   time.Duration(cfg.InputIntervalMs) * time.Millisecond,
		PackLength: cfg.PackLength,
		OutputChan: inputPacks,
		StopChan:   stopChan}

	go inputPacksGenerator.Start(cfg)
	glog.Infoln("Pack generator started")

	// Start the gRPC server in a separate goroutine
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

	// Set up and start the REST API server using Gin
	gin.SetMode(gin.ReleaseMode)
	h := rest.NewDataServiceServer(dataService)
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v1.GET("data/:id", h.GetByID)
	v1.GET("data", h.ListByTimeRange)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the REST server and listen for HTTP requests
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
