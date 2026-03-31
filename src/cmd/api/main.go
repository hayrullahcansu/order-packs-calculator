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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	db := db.InitiateSqliteDbContext()

	orderPackRepository := repository.NewOrderPackRepository(db)
	orderPackService := service.NewOrderPackService(orderPackRepository)

	// initiate handlers
	messageHandler := router.NewOrderPackHandler(orderPackService)
	srv := router.InitRouter(port, messageHandler)

	go func() {
		logging.Infof("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Infof("Server error: %v\n", err)
			cancel()
		}
	}()

	<-ctx.Done()
	// wait context cancellation and shut down resources gracefully
	logging.Warning("Stop signal received. Server is going to shut down gracefully")

	shutdownCtx, stop := context.WithTimeout(context.Background(), time.Second*30)
	defer stop()

	if stopErr := srv.Shutdown(shutdownCtx); stopErr != nil {
		logging.Errorf("error while stopping server %v\n", stopErr)
	}

	if sqlDB, err := db.DB(); err == nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			logging.Errorf("error while closing db connections %v\n", closeErr)
		}
	}

}
