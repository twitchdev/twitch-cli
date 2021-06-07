// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package mock_errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorMessage struct {
	StatusCode int    `json:"status"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

func GetErrorBytes(statusCode int, err error, message string) []byte {
	em := ErrorMessage{
		StatusCode: statusCode,
		Error:      err.Error(),
		Message:    message,
	}

	bytes, _ := json.Marshal(em)
	return bytes
}

func WriteBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(GetErrorBytes(http.StatusBadRequest, errors.New("Bad Request"), message))
}

func WriteServerError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(GetErrorBytes(http.StatusInternalServerError, errors.New("Error processing request"), message))
}
func WriteUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(GetErrorBytes(http.StatusUnauthorized, errors.New("Unauthorized"), message))
}
func WriteForbidden(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusForbidden)
	w.Write(GetErrorBytes(http.StatusForbidden, errors.New("Forbidden"), message))
}
func WriteNotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(GetErrorBytes(http.StatusNotFound, errors.New("Not Found"), message))
}
