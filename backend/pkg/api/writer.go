package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// writeJSON writes a JSON to a response writer.
//
// Also adds the Content-Length header.
func writeJSON(writer http.ResponseWriter, toBeWritten interface{}, statusCode uint) error {
	asBytes, err := json.Marshal(toBeWritten)
	if err != nil {
		return fmt.Errorf("could not marshall: %w", err)
	}

	return write(writer, asBytes, statusCode)
}

// write a []byte to a response writer.
//
// Also adds the Content-Length header.
func write(writer http.ResponseWriter, toBeWritten []byte, statusCode uint) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.Header().Set("Content-Length", strconv.Itoa(len(toBeWritten)))
	//writer.WriteHeader(int(statusCode))

	_, err := writer.Write(toBeWritten)
	if err != nil {
		return fmt.Errorf("could not write: %w", err)
	}

	return nil
}

// ErrorMessage is used to encapsulate error replies.
type ErrorMessage struct {
	Error string
}
