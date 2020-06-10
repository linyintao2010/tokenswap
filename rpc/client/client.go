package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// HTTPGet http get
func HTTPGet(url string, params, headers map[string]string, timeout int) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	addParams(req, params)
	addHeaders(req, headers)

	return doRequest(req, timeout)
}

// HTTPPost http post
func HTTPPost(url string, body interface{}, params, headers map[string]string, timeout int) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	addParams(req, params)
	addHeaders(req, headers)
	if err := addPostBody(req, body); err != nil {
		return nil, err
	}

	return doRequest(req, timeout)
}

// HTTPRawPost http raw post
func HTTPRawPost(url string, body string, params, headers map[string]string, timeout int) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	addParams(req, params)
	addHeaders(req, headers)
	if err := addRawPostBody(req, body); err != nil {
		return nil, err
	}

	return doRequest(req, timeout)
}

func addParams(req *http.Request, params map[string]string) {
	if params != nil {
		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, val := range headers {
		req.Header.Add(key, val)
	}
}

func addPostBody(req *http.Request, body interface{}) error {
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-type", "application/json")
		req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(bytes.NewBuffer(jsonData)), nil
		}
		req.Body, _ = req.GetBody()
	}
	return nil
}

func addRawPostBody(req *http.Request, body string) error {
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(body)), nil
		}
		req.Body, _ = req.GetBody()
	}
	return nil
}

func doRequest(req *http.Request, timeoutSeconds int) (*http.Response, error) {
	if timeoutSeconds <= 0 {
		return http.DefaultClient.Do(req)
	}
	timeout := time.Duration(timeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return http.DefaultClient.Do(req.WithContext(ctx))
}
