package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

/*
{
      "ip_prefix": "54.239.1.96/28",
      "region": "eu-north-1",
      "service": "AMAZON",
      "network_border_group": "eu-north-1"
	},

*/

// IPRange register
type IPRange struct {
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

// IPRanges struct
type IPRanges struct {
	SyncToken  string    `json:"syncToken"`
	CreateDate string    `json:"createDate"`
	Prefixes   []IPRange `json:"prefixes"`
}

func fetchAWSRanges() (*IPRanges, error) {
	// aws ip ranges
	// https://ip-ranges.amazonaws.com/ip-ranges.json

	var body []byte
	var err error

	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, "https://ip-ranges.amazonaws.com/ip-ranges.json", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	IPR := IPRanges{}
	if err = json.Unmarshal(body, &IPR); err != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), err.Error())
		return nil, err
	}

	return &IPR, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func main() {
	var wg sync.WaitGroup

	ipr, err := fetchAWSRanges()
	if err != nil {
		log.Fatal(err)
	}

	for _, ipp := range ipr.Prefixes {
		wg.Add(1)
		go func(iprange *IPRange) {
			defer wg.Done()
			ip, ipnet, err := net.ParseCIDR(iprange.IPPrefix)
			if err != nil {
				log.Fatal(err)
			}
			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
				fmt.Println(ip)
			}
		}(&ipp)
	}
	wg.Wait()
}
