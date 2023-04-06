package gstool

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func GetIpLocalOnce() (ip string, loc string, err error) {
	ip, loc, err = getIpLocal1()
	if err == nil {
		return
	}
	ip, loc, err = getIpLocal2()
	if err == nil {
		return
	}
	ip, loc, err = getIpLocal3()
	if err == nil {
		return
	}
	return "", "", errors.New("query failed")
}

func getIpLocal1() (string, string, error) {
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

func getIpLocal2() (string, string, error) {
	req, err := http.NewRequest("GET", "https://api.infoip.io", nil)
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
	return m["ip"].(string), m["country_short"].(string), nil
}

func getIpLocal3() (string, string, error) {
	req, err := http.NewRequest("GET", "https://ipinfo.io", nil)
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
	return m["ip"].(string), m["country"].(string), nil
}
