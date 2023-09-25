package gstool_query_ip

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type IpAddrInfo struct {
	IP          string
	CountryCode string

	Country string
	City    string
	Region  string
	ZipCode string
	Lat     float64
	Lon     float64
}

func getQueryIPFuncMap() map[string]func() string {
	return map[string]func() string{
		"tnx": genQueryIpFunc("http://tnx.nl/ip", func(res []byte) []byte {
			return bytes.Trim(res, "<>")
		}),

		"amazonaws": genQueryIpFunc("https://checkip.amazonaws.com", bytes.TrimSpace),
		"infoip":    genQueryIpFunc("https://api.infoip.io/ip", nil),
		"ipinfo":    genQueryIpFunc("https://ipinfo.io/ip", nil),
	}
}

func genQueryIpFunc(url string, formatFunc func(res []byte) []byte) func() string {
	return func() string {
		bs, err := HttpGetContentByUserAgent(url, "curl")
		if err != nil {
			return ""
		}
		if formatFunc != nil {
			bs = formatFunc(bs)
		}
		return string(bs)
	}

}
func QueryMyIP() (string, error) {
	for _, f := range getQueryIPFuncMap() {
		ip := f()
		if len(ip) > 0 {
			return ip, nil
		}
	}

	return "", errors.New("query my ip failed")
}

func getQueryIpLocFuncMap() map[string]func() *IpAddrInfo {
	return map[string]func() *IpAddrInfo{
		"infoip": func() *IpAddrInfo {
			bs, err := HttpGetContentByUserAgent("https://api.infoip.io", "curl")
			if err != nil {
				return nil
			}

			var info struct {
				IP           string  `json:"ip"`
				CountryShort string  `json:"country_short"`
				CountryLong  string  `json:"country_long"`
				City         string  `json:"city"`
				Region       string  `json:"region"`
				PostalCode   string  `json:"postal_code"`
				Latitude     float64 `json:"latitude"`
				Longitude    float64 `json:"longitude"`
			}

			if json.Unmarshal(bs, &info) != nil {
				return nil
			}

			return &IpAddrInfo{
				IP:          info.IP,
				CountryCode: info.CountryShort,
				Country:     info.CountryLong,
				City:        info.City,
				Region:      info.Region,
				ZipCode:     info.PostalCode,
				Lat:         info.Latitude,
				Lon:         info.Longitude,
			}
		},
	}
}

func QueryMyIpLoc() (*IpAddrInfo, error) {
	for _, f := range getQueryIpLocFuncMap() {
		loc := f()
		if loc != nil {
			return loc, nil
		}
	}
	return nil, errors.New("query my loc failed")
}

func HttpGetContentByUserAgent(url, userAgent string) ([]byte, error) {
	return HttpGetContentFunc(url, func(req *http.Request) {
		req.Header.Set("User-Agent", userAgent)
	})
}
