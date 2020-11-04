package controllers

import (
	"fmt"
	"imagine2/config"
	"imagine2/http"
	"runtime"

	"github.com/valyala/fasthttp"
)

// StatsController - show application stats
func StatsController(ctx *fasthttp.RequestCtx) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := map[string]interface{}{
		"memory": map[string]interface{}{
			"alloc":       formatBytes(m.Alloc),
			"total_alloc": formatBytes(m.TotalAlloc),
			"sys":         formatBytes(m.Sys),
			"num_gc":      formatBytes(uint64(m.NumGC)),
		},
		"system": m,
		"config": config.Context,
	}

	http.JSON(ctx, http.JSONResponse{
		Success: true,
		Stats:   stats,
	}, fasthttp.StatusOK)
}

func formatBytes(b uint64) string {
	value := float64(b) / 1024 / 1024
	result := fmt.Sprintf("%f", value)
	return result
}
