package gstool_query_ip

import "testing"

func TestQueryLocalIpLocInfo(t *testing.T) {
	ipinfo, err := QueryMyIpLoc()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipinfo)
}
func TestQueryLocalIpInfo(t *testing.T) {
	ipinfo, err := QueryMyIP()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ipinfo)
}
