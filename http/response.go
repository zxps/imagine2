package http

import (
	"imagine2/files"
	"imagine2/image"
	"imagine2/models"
	"imagine2/utils"
	"io"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/valyala/fasthttp"
)

// Response ...
func Response(ctx *fasthttp.RequestCtx, content []byte, code int) {
	ctx.Response.SetStatusCode(code)

	if len(content) > 0 {
		ctx.Write(content)
	}
}

// ShowFileResponse ...
func ShowFileResponse(ctx *fasthttp.RequestCtx, file *models.File, code int) {
	p := files.GetPartitionFromFile(file)

	ShowFilePartitionResponse(ctx, p, code)
}

// ShowTransformedFileResponse ...
func ShowTransformedFileResponse(ctx *fasthttp.RequestCtx, transform *files.FileTransform) {
	err := files.SynchronizePathFromFilename(transform.Target.GetFilepath())
	if err != nil {
		Response(ctx, []byte(err.Error()), fasthttp.StatusBadRequest)
		return
	}

	err = image.TransformFile(transform)
	if err != nil {
		Response(ctx, []byte(err.Error()), fasthttp.StatusBadRequest)
		return
	}

	ShowFilePartitionResponse(ctx, transform.Target, fasthttp.StatusOK)
}

// ShowFilePartitionResponse ...
func ShowFilePartitionResponse(ctx *fasthttp.RequestCtx, p *files.FilePartition, code int) {
	ctx.Response.SetStatusCode(code)

	filepath := p.GetFilepath()
	if !utils.IsFileExists(filepath) {
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	fileHandle, err := os.Open(filepath)
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	defer fileHandle.Close()

	if len(p.Mime) > 0 {
		ctx.Response.Header.SetCanonical([]byte("Content-type"), []byte(p.Mime))
	} else {
		mime, err := mimetype.DetectFile(filepath)
		if err == nil {
			ctx.Response.Header.SetCanonical([]byte("Content-type"), []byte(mime.String()))
		}
	}

	io.Copy(ctx.Response.BodyWriter(), fileHandle)
}
