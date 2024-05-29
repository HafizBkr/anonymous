package utils

import (
	"anonymous/types"
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Nature          string `json:"nature"`
	ValidationError struct {
		Field string `json:"field"`
		Error string `json:"error"`
	} `json:"validation_error"`
	ServiceError string `json:"service_error"`
}

func HandleBodyDecodingErr(w http.ResponseWriter, err error, logger types.Logger) {
	logger.Error(fmt.Sprintf(
		"Error while decoding request body %s",
		err,
	))
	WriteServiceError(w, ErrDecodingBody)
}

func WriteValidationError(w http.ResponseWriter, err map[string]string) {
	error := Error{
		Nature: "validation_error",
	}
	for k, v := range err {
		error.ValidationError.Field = k
		error.ValidationError.Error = v
		break
	}
	bytes, _ := json.Marshal(error)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bytes)
}

func WriteServiceError(w http.ResponseWriter, err error) {
	serviceErr := err.(types.ServiceError)
	error := Error{
		Nature:       "service_error",
		ServiceError: serviceErr.ErrorCode,
	}
	bytes, _ := json.Marshal(error)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(serviceErr.StatusCode)
	w.Write(bytes)
}

func WriteData(w http.ResponseWriter, statusCode int, data interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	bytes, _ := json.Marshal(map[string]interface{}{
		"data": data,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func WriteError(w http.ResponseWriter, err error) {
    serviceErr := err.(types.ServiceError)
    bytes, _ := json.Marshal(map[string]interface{}{
        "error": map[string]interface{}{
            "service":    serviceErr.ErrorCode,
            "validation": map[string]interface{}{},
        },
    })
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(serviceErr.StatusCode)
    w.Write(bytes)
}
