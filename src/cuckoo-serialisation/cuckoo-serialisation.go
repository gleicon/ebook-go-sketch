package main

import (
	"fmt"

	cuckoo "github.com/seiflotfy/cuckoofilter"
)

var safeBrowsingList *cuckoo.Filter

func testAndReport(filter *cuckoo.Filter, url string) {
	uu := []byte(url)
	if filter.Lookup(uu) {
		fmt.Println(url, "is not safe")
	} else {
		fmt.Println(url, "seems safe")
	}
}

func main() {
	safeBrowsingList = cuckoo.NewFilter(1000)
	safeBrowsingList.InsertUnique([]byte("https://badsite.com"))
	safeBrowsingList.InsertUnique([]byte("https://anotherbadsite.com"))

	testAndReport(safeBrowsingList, "https://badsite.com")
	testAndReport(safeBrowsingList, "https://anotherbadsite.com")
	testAndReport(safeBrowsingList, "https://lerolero.com")

	count := safeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

	// Delete a string (and it a miss)
	safeBrowsingList.Delete([]byte("hello"))

	count = safeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

	fmt.Println("Encoding")

	serFilter := safeBrowsingList.Encode()

	fmt.Printf("Serialized: % x\n", serFilter)

	BackupsafeBrowsingList, _ := cuckoo.Decode(serFilter)

	count = BackupsafeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

	testAndReport(BackupsafeBrowsingList, "https://badsite.com")
	testAndReport(BackupsafeBrowsingList, "https://anotherbadsite.com")
	testAndReport(BackupsafeBrowsingList, "https://lerolero.com")
}
