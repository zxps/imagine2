package controllers

import (
	"fmt"
	"imagine2/files"
	"imagine2/http"
	"imagine2/image"
	"imagine2/utils"

	"github.com/sirupsen/logrus"
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
	logrus.Info(transform.Target)
	if utils.IsFileExists(transform.Target.GetFilepath()) {
		http.ShowFilePartitionResponse(ctx, transform.Target, fasthttp.StatusOK)
	} else {
		err := files.SynchronizePathFromFilename(transform.Target.GetFilepath())
		if err != nil {
			http.Response(ctx, []byte(err.Error()), fasthttp.StatusBadRequest)
		}

		err = image.TransformFile(transform)
		if err != nil {
			http.Response(ctx, []byte(err.Error()), fasthttp.StatusBadRequest)
			return
		}

		http.ShowFilePartitionResponse(ctx, transform.Target, fasthttp.StatusOK)
	}
}
