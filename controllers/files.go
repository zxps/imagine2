package controllers

import (
	"imagine2/http"
	"imagine2/storage"

	"github.com/valyala/fasthttp"
)

// FilesController ...
func FilesController(ctx *fasthttp.RequestCtx) {
	files := storage.GetFiles(0, 50, true)

	http.JSON(ctx, http.JSONResponse{
		Success: true,
		Files:   files,
	}, fasthttp.StatusOK)
}
