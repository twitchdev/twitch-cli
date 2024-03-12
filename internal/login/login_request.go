// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package login

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/twitchdev/twitch-cli/internal/request"
)

type loginRequestResponse struct {
	StatusCode int
	Body       []byte
}

type loginHeader struct {
	Key   string
	Value string
}

func loginRequest(method string, url string, payload io.Reader) (loginRequestResponse, error) {
	return loginRequestWithHeaders(method, url, payload, []loginHeader{})
}

func loginRequestWithHeaders(method string, url string, payload io.Reader, headers []loginHeader) (loginRequestResponse, error) {
	req, err := request.NewRequest(method, url, payload)

	if err != nil {
		return loginRequestResponse{}, err
	}

	for _, header := range headers {
		req.Header.Add(header.Key, header.Value)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return loginRequestResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return loginRequestResponse{}, err
	}

	return loginRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}

func dcfInitiateRequest(url string, clientId string, scopes string) (loginRequestResponse, error) {
	formData := map[string]string{
		"client_id": clientId,
		"scopes":    scopes,
	}

	return sendMultipartPostRequest(url, formData)
}

func dcfTokenRequest(url string, clientId string, scopes string, deviceCode string, grantType string) (loginRequestResponse, error) {
	formData := map[string]string{
		"client_id":   clientId,
		"scopes":      scopes,
		"device_code": deviceCode,
		"grant_type":  grantType,
	}

	return sendMultipartPostRequest(url, formData)
}

// Creates and sends a request with the content type multipart/form-data
func sendMultipartPostRequest(url string, formData map[string]string) (loginRequestResponse, error) {
	// Create form's body using the provided data
	formBody := new(bytes.Buffer)
	mp := multipart.NewWriter(formBody)
	for k, v := range formData {
		mp.WriteField(k, v)
	}
	mp.Close() // If you do defer on this instead, it gets an "unexpected EOF" error from Twitch's servers

	req, err := request.NewRequest("POST", url, formBody)
	if err != nil {
		return loginRequestResponse{}, err
	}

	// Add Content-Type header, generated with the boundary associated with the form
	req.Header.Add("Content-Type", mp.FormDataContentType())

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return loginRequestResponse{}, err
	}

	responseBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return loginRequestResponse{}, err
	}

	return loginRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       responseBody,
	}, nil
}
