package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/axiomhq/topk"
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
	var ipCount uint64

	ipr, err := fetchAWSRanges()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Ranges: %d\n", len(ipr.Prefixes))

	tk := topk.New(100)

	ipCount = 0
	rangeCount := 0
	exactCount := make(map[string]int)

	for _, ipp := range ipr.Prefixes {
		rangeCount++
		//Uncomment to test a shorter run
		if rangeCount > 100 {
			break
		}

		ip, ipnet, err := net.ParseCIDR(ipp.IPPrefix)
		if err != nil {
			log.Println(err)
		} else {
			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
				// calculate most common first 2 octets from ip address (xx.yy)
				octetPrefix := strings.Split(ip.String(), ".")
				op := fmt.Sprintf("%s.%s", octetPrefix[0], octetPrefix[1])
				exactCount[op]++

				e := tk.Insert(op, 1)
				if e.Count < exactCount[op] {
					fmt.Printf("Error: estimate lower than exact: key=%v, exact=%v, estimate=%v\n", e.Key, exactCount[op], e.Count)
				}
				if e.Count-e.Error > exactCount[op] {
					fmt.Printf("Error: error bounds too large: key=%v, count=%v, error=%v, exact=%v\n", e.Key, e.Count, e.Error, exactCount[op])
				}
			}
		}

	}
	fmt.Printf("Unique IP Addresses: %d\n", ipCount)
	el := tk.Estimate("52.94")
	fmt.Printf("Prefix 52.94 ranks at %d\n ", el.Count)

	kk := tk.Keys()
	fmt.Println("List top 10 matches\nFirst 2 IP octets - Count")
	for i := 0; i < 10; i++ {
		fmt.Printf("%s - %d\n", kk[i].Key, kk[i].Count)
	}

}
