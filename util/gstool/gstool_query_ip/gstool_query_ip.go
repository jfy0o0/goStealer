package gstool_query_ip

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type IpInfo struct {
	IP          string
	CountryCode string

	Country string
	City    string
	Region  string
	ZipCode string
	Lat     float64
	Lon     float64
}

func QueryIpInfo(ip string) *IpInfo {

	funcs := make([]func(string) *IpInfo, 0, 2)

	if (time.Now().Nanosecond() & 0x10) == 0 {
		funcs = append(funcs, queryByApiNinjas, queryByInfoIP)
	} else {
		funcs = append(funcs, queryByInfoIP, queryByApiNinjas)
	}

	for _, f := range funcs {
		info := f(ip)
		if info != nil {
			return info
		}
	}

	return nil
}

func queryByApiNinjas(ip string) *IpInfo {
	res, err := httpGetFunc("https://api.api-ninjas.com/v1/iplookup?address="+ip, func(req *http.Request) {
		req.Header.Set("X-Api-Key", "TP191s9ORkohzwQVs1koHA==AoMCf7DTOot4pj4y")
	})
	if err != nil {
		return nil
	}

	m := make(map[string]interface{}, 10)

	if err = json.Unmarshal(res, &m); err != nil {
		return nil
	}

	if _, ok := m["is_valid"]; !ok {
		return nil
	}

	return &IpInfo{
		IP:          ip,
		CountryCode: m["country_code"].(string),
		Country:     m["country"].(string),
		City:        m["city"].(string),
		Region:      m["region"].(string),
		ZipCode:     m["zip"].(string),
		Lat:         m["lat"].(float64),
		Lon:         m["lon"].(float64),
	}
}

func queryByInfoIP(ip string) *IpInfo {
	res, err := httpGetFunc("https://api.infoip.io/"+ip, nil)
	if err != nil {
		return nil
	}
	m := make(map[string]interface{}, 10)

	if err = json.Unmarshal(res, &m); err != nil {
		return nil
	}

	return &IpInfo{
		IP:          ip,
		CountryCode: m["country_short"].(string),
		Country:     m["country_long"].(string),
		City:        m["city"].(string),
		Region:      m["region"].(string),
		ZipCode:     m["postal_code"].(string),
		Lat:         m["latitude"].(float64),
		Lon:         m["longitude"].(float64),
	}
}

func httpGetFunc(url string, formatReqFunc func(*http.Request)) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	userAgent := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/94.0.835.163 Safari/535.1"
	req.Header.Set("User-Agent", userAgent)

	if formatReqFunc != nil {
		formatReqFunc(req)
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	if resp, e := client.Do(req); e != nil {
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
