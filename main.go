package ipinfo

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
)

type IP struct {
	Ip            string `json:"ip"`
	City          string `json:"city"`
	Organization  string `json:"organization"`
	Asn           string `json:"asn"`
	Region        string `json:"region"`
	RegionCode    string `json:"region_code"`
	Country       string `json:"country"`
	CountryCode   string `json:"country_code"`
	CountryCode3  string `json:"country_code3"`
	ContinentCode string `json:"continent_code"`
	PostalCode    string `json:"postal_code"`
	Timezone      string `json:"timezone"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	IsPrivate     bool   `json:"isPrivate,omitempty"`
	Str           string `json:"str,omitempty"`
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

	apiUrl := fmt.Sprintf("https://api.seeip.org/geoip/%s", ipStr)
	if slices.Contains([]string{"ipv4", "ipv6"}, ipStr) {
		apiUrl = fmt.Sprintf("https://%s.seeip.org/geoip/", ipStr)
	}

	resp, _ := http.Get(apiUrl)
	data := IP{Ip: ipStr}
	_ = json.NewDecoder(resp.Body).Decode(&data)
	data.Str = fmt.Sprintf("%s - %s - %s", data.Ip, data.Organization, strings.Join([]string{data.City, data.Region, data.Country}, ", "))
	return data
}

func ParseMyIP() IP {
	return ParseIP("")
}

func ParseMyIPv4() IP {
	return ParseIP("ipv4")
}

func ParseMyIPv6() IP {
	return ParseIP("ipv6")
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
