package opsysinfo

import (
	"encoding/json"
	"io"
	"net/http"
)

type IpInfo struct {
	IP          string `json:"ip"`
	Country     string `json:"country_name"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Org         string `json:"org"`
}

type MacInfo struct {
	Company string `json:"company"`
}

const (
	ipAPI     string = "https://ipapi.co/json"
	macAPI    string = "https://www.macvendorlookup.com/api/v2"
	userAgent string = "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.7113.93 Safari/537.36"
)

// Public IP API Lookup
func PublicIPInfo() IpInfo {
	var data = IpInfo{}
	req, err := http.NewRequest(http.MethodGet, ipAPI, nil)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		return data
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return data
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	return data
}

// Mac Address API lookup
func MacAddressVendor(macAddress string) MacInfo {
	var data = MacInfo{}
	req, err := http.NewRequest(http.MethodGet, macAPI+macAddress, nil)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		return data
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return data
	}
	defer resp.Body.Close()

	macVendorBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(macVendorBody, &data)

	return data
}

// Show country flag
func CountryFlag() string {

	ipInfo := PublicIPInfo()

	for k, v := range flag {
		if ipInfo.CountryCode == k {
			return v
		}
	}

	return flag["Unknown"]
}
