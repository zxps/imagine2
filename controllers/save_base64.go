package controllers

import (
	"imagine2/files"
	"imagine2/http"
	"imagine2/image"
	"imagine2/models"
	"imagine2/storage"
	"imagine2/utils"
	"strconv"

	"github.com/valyala/fasthttp"
)

// SaveBase64 - save(upload) or update a file in storage with base64 data
func SaveBase64(ctx *fasthttp.RequestCtx) {
	fileIDParam := string(ctx.FormValue("id"))

	var err error
	var file *models.File

	if len(fileIDParam) > 0 {
		fileID, err := strconv.ParseInt(fileIDParam, 10, 0)
		if err != nil || fileID < 1 {
			http.JSONStatus(ctx, "bad file id", fasthttp.StatusBadRequest)
			return
		}

		file, err = storage.GetFileByID(int(fileID))
		if err != nil {
			http.JSONStatus(ctx, err.Error(), fasthttp.StatusBadRequest)
			return
		}

		data := string(ctx.FormValue("data"))
		partition := files.GetPartitionFromFile(file)
		files.UploadFileFromBase64(partition, data)

		image.ProcessImagePartition(partition, image.ProcessImageOptions{
			Normalize: true,
		})

		file.Size = partition.Size
		file.Width = partition.ImageWidth
		file.Height = partition.ImageHeight
		file.Updated = utils.StorageTimestamp()

	} else {
		partition := files.NewFilePartition(string(ctx.FormValue("name")))
		partition.Generate()
		if !partition.Success {
			http.JSONStatus(ctx, "Unable to generate new file partition", fasthttp.StatusBadRequest)
			return
		}

		err := partition.Sync()
		if err != nil {
			http.JSONStatus(ctx, "Unable to synchronize file partition", fasthttp.StatusBadRequest)
			return
		}

		err = files.UploadFileFromBase64(partition, string(ctx.FormValue("data")))
		if err != nil {
			http.JSONStatus(ctx, err.Error(), fasthttp.StatusBadRequest)
			return
		}

		image.ProcessImagePartition(partition, image.ProcessImageOptions{
			Normalize: true,
		})

		file = storage.TransformPartitionToFile(partition)
	}

	err = storage.SaveFile(file)

	code := fasthttp.StatusOK
	status := ""
	if err != nil {
		status = err.Error()
		code = fasthttp.StatusBadRequest
	} else {
		files.InvalidateFileTransforms(*file)
	}

	http.JSONFile(ctx, file, status, code)
}
