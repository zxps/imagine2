package controllers

import (
	"imagine2/files"
	"imagine2/http"
	"imagine2/image"
	"imagine2/storage"
	"imagine2/tasks"

	"github.com/valyala/fasthttp"
)

// SaveBase64 - save(upload) or update a file in storage with base64 data
func SaveBase64(ctx *fasthttp.RequestCtx) {
	p, err := files.UploadFileFromBase64(
		string(ctx.FormValue("data")),
		string(ctx.FormValue("name")),
	)

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

	tasks.NotifyFileCreated(*file)

	http.JSONFile(ctx, file, status, code)
}
