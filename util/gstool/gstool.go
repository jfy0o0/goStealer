package gstool

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetIpLocal() (string, string, error) {
	req, err := http.NewRequest("GET", "http://ip-api.com/json?lang=zh-CN", nil)
	if err != nil {
		return "", "", err
	}
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return "", "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()
		return "", "", err
	}
	res.Body.Close()
	m := make(map[string]interface{}, 100)
	if err = json.Unmarshal(body, &m); err != nil {
		return "", "", err
	}
	if m["status"] != "success" {
		return "", "", err
	}
	ip := m["query"]
	countryCode := m["countryCode"]

	return ip.(string), countryCode.(string), nil
}
