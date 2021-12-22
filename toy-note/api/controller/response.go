package controller

import "fmt"

type errorMessage struct {
	Err string `json:"error"`
}

type successMessage struct {
	Success string `json:"success"`
}

func errorResponse(err error) errorMessage {
	return errorMessage{Err: err.Error()}
}

func successResponse(data interface{}) successMessage {
	return successMessage{Success: fmt.Sprintf("%v", data)}
}

type downloadSuccess struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}
