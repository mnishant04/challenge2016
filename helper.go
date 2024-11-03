package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func SendResponse(statusCode int, message string, data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(statusCode)
	apiResp := Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	responseBytes, err := json.Marshal(apiResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseBytes)
}

func UnmarshalJson[T any](body io.Reader) (val T, err error) {
	data, err := io.ReadAll(body)
	if err != nil {
		log.Printf("Error while reading body: %v", err)
		return val, err
	}
	err = json.Unmarshal(data, &val)
	if err != nil {
		log.Printf("error while unmarshalling: %s", err)
		return val, err
	}
	return
}
