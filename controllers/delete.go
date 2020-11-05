package controllers

import (
	"imagine2/files"
	"imagine2/http"
	"imagine2/storage"
	"imagine2/tasks"
	"os"
	"strconv"

	"github.com/valyala/fasthttp"
)

// Delete - Delete file controller
func Delete(ctx *fasthttp.RequestCtx) {
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

	partition := files.GetPartitionFromFile(file)

	err = os.Remove(partition.GetFilepath())
	if err != nil {
		http.JSONError(ctx, err)
		return
	}

	err = storage.DeleteFile(file.ID)
	if err != nil {
		http.JSONError(ctx, err)
		return
	}

	transforms, _ := files.GetAllFileTransforms(*file)

	for _, t := range transforms {
		os.Remove(t.FullFilepath)
	}

	tasks.NotifyFileDelete(*file)

	http.JSON(ctx, http.JSONResponse{Success: true}, fasthttp.StatusOK)
}
