// Package router configures HTTP routing, static file serving, and API endpoint registration.
package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

// resolveStaticDir determines the static assets directory.
// In production (Docker), it uses ./static next to the binary.
// In development (go run), it falls back to the source tree path.
func resolveStaticDir() string {
	// prefer ./static next to the binary (Docker / production)
	if _, err := os.Stat("static/index.html"); err == nil {
		return "static"
	}
	// fallback: resolve relative to source file (go run / development)
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), "..", "static")
}

// InitRouter creates and configures the HTTP server with all routes and middleware.
func InitRouter(
	port *int,
	orderPackHandler *OrderPackHandler,
) *http.Server {
	r := gin.Default()
	// middleware recovers from any panic and returns 500 error code.
	r.Use(gin.Recovery())

	// serve static frontend
	staticDir := resolveStaticDir()
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
