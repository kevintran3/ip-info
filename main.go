package ipinfo

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type IP struct {
	Ip          string `json:"query"`
	City        string `json:"city"`
	Isp         string `json:"isp"`
	Org         string `json:"org"`
	As          string `json:"as"`
	RegionName  string `json:"regionName"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Zip         string `json:"zip"`
	Timezone    string `json:"timezone"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	IsPrivate   bool   `json:"isPrivate,omitempty"`
	Str         string `json:"str,omitempty"`
}

func ParseIP(ipStr string) IP {
	if len(ipStr) > 0 {
		ip := net.ParseIP(ipStr)
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return IP{
				Ip:        ipStr,
				IsPrivate: true,
				Str:       fmt.Sprintf("%s - private IP", ipStr),
			}
		}
	}
	resp, _ := http.Get("http://ip-api.com/json/" + ipStr)
	data := IP{Ip: ipStr}
	_ = json.NewDecoder(resp.Body).Decode(&data)
	data.Str = fmt.Sprintf("%s - %s - %s", data.Ip, data.Isp, strings.Join([]string{data.City, data.RegionName, data.Country}, ", "))
	return data
}

func ParseMyIP() IP {
	return ParseIP("")
}

func ClientIP(req *http.Request) string {
	ipHeaders := []string{
		"Cf-Connecting-Ip",
		"Fastly-Client-Ip",
		"X-Appengine-User-Ip",
		"X-Real-Ip",
		"X-Forwarded-For",
	}
	for _, h := range ipHeaders {
		if ip, ok := req.Header[h]; ok {
			return ip[0]
		}
	}
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	return ip
}

func ParseClientIP(req *http.Request) IP {
	return ParseIP(ClientIP(req))
}
