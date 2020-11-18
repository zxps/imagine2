package controllers

import (
	"fmt"
	"imagine2/files"
	"imagine2/http"
	"imagine2/utils"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// Render - render file by path
func Render(ctx *fasthttp.RequestCtx) {
	filepath := fmt.Sprintf("%v", ctx.UserValue("filepath"))

	logrus.Info("filepath is:", filepath)

	if len(filepath) < 1 {
		http.Response(ctx, []byte(""), fasthttp.StatusBadRequest)
		return
	}

	transform := files.ExtractTransform(filepath)

	if utils.IsFileExists(transform.Target.GetFilepath()) {
		http.ShowFilePartitionResponse(ctx, transform.Target, fasthttp.StatusOK)
		return
	}

	if !utils.IsFileExists(transform.Source.GetFilepath()) {
		http.Response(ctx, []byte("file not found"), fasthttp.StatusNotFound)
		return
	}

	http.ShowTransformedFileResponse(ctx, transform)
}
