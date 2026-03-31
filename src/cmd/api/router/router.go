package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	port *int,
	orderPackHandler *OrderPackHandler,
) *http.Server {
	r := gin.Default()
	// middleware recovers from any panic and returns 500 error code.
	r.Use(gin.Recovery())

	// serve static frontend
	_, currentFile, _, _ := runtime.Caller(0)
	staticDir := filepath.Join(filepath.Dir(currentFile), "..", "static")
	r.StaticFile("/", staticDir+"/index.html")
	r.Static("/static", staticDir)

	v1 := r.Group("v1")
	{
		v1OrderPacks := v1.Group("order_packs")
		{
			v1OrderPacks.GET("", orderPackHandler.GetAvailableOrderPacks)
			v1OrderPacks.POST("", orderPackHandler.AddOrderPack)
			v1OrderPacks.PUT("/:id", orderPackHandler.UpdateOrderPack)
			v1OrderPacks.DELETE("/:id", orderPackHandler.RemoveOrderPack)
			v1OrderPacks.POST("solve", orderPackHandler.SolveOrderPacks)
		}
	}

	url := fmt.Sprintf("0.0.0.0:%d", *port)
	// Create HTTP server
	srv := &http.Server{
		Addr:    url,
		Handler: r,
	}

	return srv
}
