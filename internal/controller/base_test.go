package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

func tPreparePostRequest(url string, body any) *http.Request {
	jsonBody, _ := json.Marshal(body)

	return httptest.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
}

func tPrepareGetRequest(url string) *http.Request {
	return httptest.NewRequest(http.MethodGet, url, nil)
}

func tPrepareResponse(buffer *bytes.Buffer, response any) {
	responseBytes, _ := io.ReadAll(buffer)

	_ = json.Unmarshal(responseBytes, &response)
}
