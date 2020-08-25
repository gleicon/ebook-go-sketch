package main

import (
	"fmt"

	"github.com/yourbasic/bloom"
)

var safeBrowsingList *bloom.Filter

func testAndReport(url string) {
	if safeBrowsingList.Test(url) {
		fmt.Println(url, "is not safe")
	} else {
		fmt.Println(url, "seems safe")
	}
}
func main() {
	// 1000 elementos, erro de 1/20 (0.5%)
	safeBrowsingList = bloom.New(1000, 20)

	safeBrowsingList.Add("https://badsite.com")
	safeBrowsingList.Add("https://anotherbadsite.com")

	fmt.Printf("Sites: %d\n", safeBrowsingList.Count())

	testAndReport("https://lerolero.com")
	testAndReport("https://badsite.com")

}
