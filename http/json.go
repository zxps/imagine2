package http

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// JSONResponse - standard response
type JSONResponse struct {
	Success bool        `json:"success"`
	Status  string      `json:"status,omitempty"`
	File    interface{} `json:"file,omitempty"`
	Files   interface{} `json:"files,omitempty"`
	Stats   interface{} `json:"stats,omitempty"`
}

// AllowCORS - allow all CORS
func AllowCORS(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))
}

// JSON ...
func JSON(ctx *fasthttp.RequestCtx, response JSONResponse, code int) {
	writeJSONHeaders(ctx, code)
	writeJSON(ctx, response)
}

// JSONSuccess - send success json response
func JSONSuccess(ctx *fasthttp.RequestCtx) {
	writeJSONHeaders(ctx, fasthttp.StatusOK)
	writeJSON(ctx, &JSONResponse{
		Success: isSuccessResponse(fasthttp.StatusOK),
		Status:  "",
	})
}

// JSONStatus - send json response
func JSONStatus(ctx *fasthttp.RequestCtx, status string, code int) {
	writeJSONHeaders(ctx, code)
	writeJSON(ctx, &JSONResponse{
		Success: isSuccessResponse(code),
		Status:  status,
	})
}

// JSONError - send error
func JSONError(ctx *fasthttp.RequestCtx, e error) {
	writeJSONHeaders(ctx, fasthttp.StatusBadRequest)
	writeJSON(ctx, &JSONResponse{
		Success: isSuccessResponse(fasthttp.StatusBadRequest),
		Status:  e.Error(),
	})
}

// JSONStats - send json response from interface structure
func JSONStats(ctx *fasthttp.RequestCtx, stats interface{}, status string, code int) {
	writeJSONHeaders(ctx, code)

	response := &JSONResponse{
		Success: isSuccessResponse(code),
	}

	if len(status) > 0 {
		response.Status = status
	}

	if stats != nil {
		response.Stats = stats
	}

	writeJSON(ctx, response)
}

// JSONFile - send json response from interface structure
func JSONFile(ctx *fasthttp.RequestCtx, file interface{}, status string, code int) {
	writeJSONHeaders(ctx, code)

	response := &JSONResponse{
		Success: isSuccessResponse(code),
	}

	if len(status) > 0 {
		response.Status = status
	}

	if file != nil {
		response.File = file
	}

	writeJSON(ctx, response)
}

func writeJSONHeaders(ctx *fasthttp.RequestCtx, code int) {
	ctx.Response.Header.SetCanonical([]byte("Content-type"), []byte("application/json"))
	ctx.Response.SetStatusCode(code)
}

func writeJSON(ctx *fasthttp.RequestCtx, v interface{}) {
	resultJSON, err := json.Marshal(v)
	if err != nil {
		ctx.WriteString("")
		log.Error("json marshaling error", err.Error())
		return
	}

	ctx.Write(resultJSON)
}

func isSuccessResponse(code int) bool {
	isSuccess := true
	if code < 200 || code > 299 {
		isSuccess = false
	}

	return isSuccess
}
