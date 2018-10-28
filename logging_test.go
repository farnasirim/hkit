package hkit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
	"strings"
	"testing"
)

func TestLoggerSimple(t *testing.T) {

	expRequestMap := map[string]interface{}{
		"hello": "wow",
		"complex": []interface{}{
			map[string]interface{}{
				"1": "2",
			},
			"wow",
		},
	}

	expResponseMap := map[string]interface{}{
		"key": "value",
		"another key": []interface{}{
			"val", "ues",
		},
	}
	expRequestString, _ := json.Marshal(expRequestMap)

	testHandler := func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(expResponseMap)
	}

	logWriter := bytes.NewBuffer(nil)

	handler := NewLogger(testHandler)
	handler.SetWriter(logWriter)

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	url, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	target := "/some/target/url"
	url.Path = path.Join(url.Path, target)

	http.Post(url.String(), "application/json", bytes.NewBuffer(expRequestString))

	expectedJSONs := []interface{}{expRequestMap, expResponseMap}
	actualJSONs := []interface{}{}

	for _, line := range strings.Split(logWriter.String(), "\n") {
		mp := make(map[string]interface{})
		err := json.Unmarshal([]byte(line), &mp)
		if err != nil {
			continue
		}
		actualJSONs = append(actualJSONs, mp)
	}

	if !reflect.DeepEqual(expectedJSONs, actualJSONs) {
		t.Error("expected logged output not equal to the actual one")
	}
}
