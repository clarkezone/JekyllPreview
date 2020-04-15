package main

import (
	"testing"
)

func TestInitShareManager(t *testing.T) {
	sm := createShareManager()

	sm.start()
	sm.shareBranch("/master/", "./test/one")
	sm.shareBranch("/pepper/", "./test/two")

	// c := http.Client{}
	// resp, err := c.Get("http://localhost:8085/master")
	// if err != nil {
	// 	t.Fatalf("Request failed")
	// }

	// bytes, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Fatalf("Request failed")
	// }

	// if len(bytes) < 40 {
	// 	t.Fatalf("Request failed")
	// }

}
