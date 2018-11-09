package hkit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
	"testing"
)

var (
	expRequestMap = map[string]interface{}{
		"hello": "wow",
		"complex": []interface{}{
			map[string]interface{}{
				"1": "2",
			},
			"wow",
		},
	}

	expResponseMap = map[string]interface{}{
		"key": "value",
		"another key": []interface{}{
			"val", "ues",
		},
	}
)

func TestHTTPResponseInterceptorSimple(t *testing.T) {
	expRequestString, _ := json.Marshal(expRequestMap)

	origCalled := 0

	originalHandler := func(w http.ResponseWriter, r *http.Request) {
		origCalled += 1
		json.NewEncoder(w).Encode(expResponseMap)
	}

	var cachedValue []byte = nil
	cachedHandler := func(w http.ResponseWriter, r *http.Request) {
		if cachedValue == nil {
			interceptedWriter := NewHTTPResponseReplayer(w)
			originalHandler(interceptedWriter, r)
			var err error
			cachedValue, err = interceptedWriter.Marshal()
			if err != nil {
				t.Error(err.Error())
			}
			return
		}
		replayer, err := ReplayerFromSerialized(cachedValue)
		if err != nil {
			t.Error(err.Error())
		}
		replayer.Replay(w)
	}

	testServer := httptest.NewServer(http.HandlerFunc(cachedHandler))
	defer testServer.Close()

	url, err := url.Parse(testServer.URL)
	if err != nil {
		t.Error(err.Error())
	}

	target := "/some/target/url"
	url.Path = path.Join(url.Path, target)

	resp1, err := http.Post(url.String(), "application/json",
		bytes.NewBuffer(expRequestString))
	if err != nil {
		t.Error(err.Error())
	}

	resp2, err := http.Post(url.String(), "application/json",
		bytes.NewBuffer(expRequestString))
	if err != nil {
		t.Error(err.Error())
	}

	if 1 != origCalled {
		t.Errorf("incorrect number of calls to the original handler, 1 != %d\n",
			origCalled)
	}

	body1, err1 := ioutil.ReadAll(resp1.Body)
	body2, err2 := ioutil.ReadAll(resp2.Body)

	if err1 != nil {
		t.Error(err1.Error())
	}

	if err2 != nil {
		t.Error(err2.Error())
	}

	if len(body1) != len(body2) {
		t.Errorf("different number of bytes read in subsequent calls\n")
	}

	if !reflect.DeepEqual(body1, body2) {
		t.Errorf("different body received after cache")
	}
}
