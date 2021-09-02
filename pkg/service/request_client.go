package service

import (
	"fmt"
	"io"
	"net/http"
)

func RequestClient(method, destination string, body io.ReadCloser) (resp *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest(method, destination, body)
	if err != nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{StatusCode: http.StatusServiceUnavailable}
	}
	return
}
