// Package main is the entry point for the Order Packs Calculator API server.
// It initializes the database, wires up dependencies, and handles graceful shutdown.
package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hayrullahcansu/order-packs-calculator/src/cmd/api/router"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/repository"
	"github.com/hayrullahcansu/order-packs-calculator/src/internal/service"
	"github.com/hayrullahcansu/order-packs-calculator/src/shared/db"
	"github.com/hayrullahcansu/order-packs-calculator/src/shared/logging"
)

func main() {
	port := flag.Int("port", 8080, "Api port")
	flag.Parse()

	// listen for OS interrupt/kill signals for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// initialize database with auto-migration and seed data
	db := db.InitiateSqliteDbContext()

	// wire up dependencies: repository -> service -> handler -> router
	orderPackRepository := repository.NewOrderPackRepository(db)
	orderPackService := service.NewOrderPackService(orderPackRepository)
	messageHandler := router.NewOrderPackHandler(orderPackService)
	srv := router.InitRouter(port, messageHandler)

	// start HTTP server in a separate goroutine
	go func() {
		logging.Infof("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Infof("Server error: %v\n", err)
			cancel()
		}
	}()

	// block until shutdown signal is received
	<-ctx.Done()
	logging.Warning("Stop signal received. Server is going to shut down gracefully")

	// allow up to 30 seconds for in-flight requests to complete
	shutdownCtx, stop := context.WithTimeout(context.Background(), time.Second*30)
	defer stop()

	if stopErr := srv.Shutdown(shutdownCtx); stopErr != nil {
		logging.Errorf("error while stopping server %v\n", stopErr)
	}

	// close the underlying database connection pool
	if sqlDB, err := db.DB(); err == nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			logging.Errorf("error while closing db connections %v\n", closeErr)
		}
	}

}
