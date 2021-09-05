package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

// DecodeReqRespBody decodes request or response arguments from request/response body to a schema
func DecodeReqRespBody(body io.Reader, v interface{}) error {
	if body == nil {
		return errors.New("error decoding: body is nil")
	}
	if v == nil {
		return errors.New("error decoding: interface is nil")
	}
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decoding:%s", err.Error())
	}
	return nil
}

// EncodeReqRespBody encodes request or response arguments from a schema to request/response body
func EncodeReqRespBody(body interface{}) (io.ReadCloser, error) {
	if body == nil {
		return nil, errors.New("error encoding: request body is nil")
	}
	b, _ := json.Marshal(body)
	return ioutil.NopCloser(bytes.NewReader(b)), nil

}
