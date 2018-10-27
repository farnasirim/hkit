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

	testHandler := func(w http.ResponseWriter, req *http.Request) {
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

	payload := map[string]interface{}{"hello": "wow", "complex": []interface{}{map[string]interface{}{"1": "2"}, 2, 3, "wow"}}
	payloadString, err := json.Marshal(payload)
	if err != nil {
		t.Error(err.Error())
	}

	payloadReader := bytes.NewBuffer(payloadString)

	_, err = http.Post(url.String(), "application/json", payloadReader)
	if err != nil {
		t.Error(err.Error())
	}

	mapFromLogJSON := func(log string) map[string]interface{} {
		parts := strings.Split(log, "\n")
		lastPart := parts[len(parts)-1]
		mp := make(map[string]interface{})
		if err := json.Unmarshal([]byte(lastPart), &mp); err != nil {
			t.Error(err.Error())
		}
		return mp
	}

	expected := `Method: POST
remote address: [::1]:57770

User-Agent: curl/7.54.1
Accept: */*
Content-Length: 16
Content-Type: application/x-www-form-urlencoded

`
	expected += string(payloadString)

	if !reflect.DeepEqual(mapFromLogJSON(expected), mapFromLogJSON(logWriter.String())) {
		t.Error("expected logged output not equal to the actual one")
	}
}
