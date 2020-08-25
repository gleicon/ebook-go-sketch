package main

import (
	"fmt"

	cuckoo "github.com/seiflotfy/cuckoofilter"
)

var safeBrowsingList *cuckoo.Filter

func testAndReport(url string) {
	uu := []byte(url)
	if safeBrowsingList.Lookup(uu) {
		fmt.Println(url, "is not safe")
	} else {
		fmt.Println(url, "seems safe")
	}
}

func main() {
	safeBrowsingList = cuckoo.NewFilter(1000)
	safeBrowsingList.InsertUnique([]byte("https://badsite.com"))
	safeBrowsingList.InsertUnique([]byte("https://anotherbadsite.com"))

	testAndReport("https://badsite.com")
	testAndReport("https://anotherbadsite.com")
	testAndReport("https://lerolero.com")

	count := safeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

	// Delete a string (and it a miss)
	safeBrowsingList.Delete([]byte("hello"))

	count = safeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

	// Delete a string (a hit)
	safeBrowsingList.Delete([]byte("https://badsite.com"))

	count = safeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

	safeBrowsingList.Reset() // reset

	count = safeBrowsingList.Count()
	fmt.Printf("Items: %d\n", count)

}
