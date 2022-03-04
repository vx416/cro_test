package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type HttpHelper struct {
}

func (tester HttpHelper) GetResponseData(record *httptest.ResponseRecorder, target interface{}) error {
	data, err := ioutil.ReadAll(record.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func (tester HttpHelper) BuildRequest(method, path string, data interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var (
		req    *http.Request
		reader io.Reader
	)
	if data != nil {
		switch val := data.(type) {
		case []byte:
			reader = bytes.NewReader(val)
		case map[string]interface{}:
			jsonData, err := json.Marshal(val)
			if err != nil {
				panic(err)
			}
			reader = bytes.NewReader(jsonData)
		default:
			jsonData, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}
			reader = bytes.NewReader(jsonData)
		}
		req = httptest.NewRequest(method, path, reader)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func (tester HttpHelper) SetBearToken(req *http.Request, token string) {
	req.Header.Set("Authentication", "Bearer "+token)
}
