package controllers

import (
	"fmt"
	"imagine2/files"
	"imagine2/http"
	"imagine2/utils"

	"github.com/valyala/fasthttp"
)

// RenderController ...
func RenderController(ctx *fasthttp.RequestCtx) {
	filepath := fmt.Sprintf("%v", ctx.UserValue("filepath"))
	if len(filepath) < 1 {
		http.Response(ctx, []byte(""), fasthttp.StatusBadRequest)
		return
	}

	transform := files.ExtractTransform(filepath)

	if utils.IsFileExists(transform.Target.GetFilepath()) {
		http.ShowFilePartitionResponse(ctx, transform.Target, fasthttp.StatusOK)
	} else {
		http.ShowTransformedFileResponse(ctx, transform)
	}
}
