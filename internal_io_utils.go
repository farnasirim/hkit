package hkit

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func readAndRecreateReadCloser(reader io.ReadCloser) ([]byte, io.ReadCloser, error) {
	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}
	newReadCloser := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, newReadCloser, nil
}

func readCloserFromJSON(jsonObj interface{}) (io.ReadCloser, error) {
	buffer, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, err
	}
	reader := ioutil.NopCloser(bytes.NewBuffer(buffer))
	return reader, nil
}
