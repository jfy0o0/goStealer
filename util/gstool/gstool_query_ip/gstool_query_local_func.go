package gstool_query_ip

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpGetContentFunc(url string, f func(*http.Request)) ([]byte, error) {

	if resp, e := HttpGetFunc(url, f); e != nil {
		return nil, e
	} else if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		if bs, e := io.ReadAll(resp.Body); e != nil {
			return nil, e
		} else {
			return bs, nil
		}
	} else {
		return nil, errors.New(resp.Status)
	}
}

func HttpGetFunc(url string, f func(*http.Request)) (*http.Response, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if f != nil {
		f(req)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client.Do(req)
}

func HttpPostFunc(url string, values url.Values, f func(*http.Request)) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if f != nil {
		f(req)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client.Do(req)
}
