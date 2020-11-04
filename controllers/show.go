package controllers

import (
	"imagine2/files"
	"imagine2/http"
	"imagine2/storage"
	"strconv"

	"github.com/valyala/fasthttp"
)

// ShowController - Show file from storage
func ShowController(ctx *fasthttp.RequestCtx) {
	fileIDParam := string(ctx.FormValue("id"))
	transformParam := string(ctx.FormValue("transform"))

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

	if len(transformParam) > 0 {
		transformFilepath := file.Path + "/" + transformParam + "/" + file.Fullname
		transform := files.ExtractTransform(transformFilepath)
		ctx.Response.Header.Set("Imagine2-Filepath", transformFilepath)
		http.ShowTransformedFileResponse(ctx, transform)
	} else {
		http.ShowFileResponse(ctx, file, fasthttp.StatusOK)
	}
}
