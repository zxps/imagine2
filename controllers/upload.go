package controllers

import (
	"imagine2/files"
	"imagine2/http"
	"imagine2/image"
	"imagine2/storage"

	"github.com/valyala/fasthttp"
)

// UploadController - upload controller
func UploadController(ctx *fasthttp.RequestCtx) {
	p, err := files.UploadFileFromRequest(ctx, "file")

	if err != nil {
		http.JSONStatus(ctx, err.Error(), fasthttp.StatusBadRequest)
		return
	}

	image.ProcessImagePartition(p, image.ProcessImageOptions{
		Normalize: true,
	})

	file, err := storage.CreateFromPartition(p)

	code := fasthttp.StatusOK
	status := ""
	if err != nil {
		status = err.Error()
		code = fasthttp.StatusBadRequest
	}

	http.JSONFile(ctx, file, status, code)
}
