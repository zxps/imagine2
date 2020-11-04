package controllers

import (
	"imagine2/files"
	"imagine2/http"
	"imagine2/storage"
	"strconv"

	"github.com/valyala/fasthttp"
)

// InvalidateCacheController ...
func InvalidateCacheController(ctx *fasthttp.RequestCtx) {
	fileIDParam := string(ctx.FormValue("id"))
	if len(fileIDParam) < 1 {
		http.JSONStatus(ctx, "no parameters", fasthttp.StatusBadRequest)
		return
	}

	fileID, err := strconv.ParseInt(fileIDParam, 10, 0)
	if err != nil || fileID < 1 {
		http.JSONStatus(ctx, "bad file id", fasthttp.StatusBadRequest)
		return
	}

	file, err := storage.GetFileByID(int(fileID))
	if err != nil {
		http.JSONStatus(ctx, err.Error(), fasthttp.StatusNotFound)
		return
	}

	invalidateCount := files.InvalidateFileTransforms(*file)

	http.JSON(ctx, http.JSONResponse{
		Success: true,
		Stats: map[string]interface{}{
			"invalidated": invalidateCount,
		},
	}, fasthttp.StatusOK)
}
