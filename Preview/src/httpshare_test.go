package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func _testInitShareManager(t *testing.T) {
	sm := createShareManager()

	sm.start()
	sm.shareBranch("/master/", "./test/one")
	//sm.shareBranch("/pepper/", "./test/two")

	c := http.Client{}
	resp, err := c.Get("http://localhost:8085/test/one/index.html")
	if err != nil {
		t.Fatalf("Request failed")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	st := string(bytes)
	fmt.Printf("%v", st)
	if err != nil {
		t.Fatalf("Request failed")
	}

	if len(bytes) < 10 {
		t.Fatalf("Request failed")
	}

}

func TestInitShareManagerRootDir(t *testing.T) {
	sm := createShareManager()

	sm.start()
	sm.shareRootDir("./test/one")

	c := http.Client{}
	resp, err := c.Get("http://localhost:8085/index.html")
	if err != nil {
		t.Fatalf("Request failed")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	st := string(bytes)
	fmt.Printf("%v", st)
	if err != nil {
		t.Fatalf("Request failed")
	}

	if len(bytes) < 10 {
		t.Fatalf("Request failed")
	}

}

// func basicAuth() string {
//     var username string = "someuser"
//     var passwd string = "somepassword"
//     client := &http.Client{}
//     req, err := http.NewRequest("GET", "http://0.0.0.0:7000", nil)
//     req.SetBasicAuth(username, passwd)
//     resp, err := client.Do(req)
//     if err != nil{
//         log.Fatal(err)
//     }
//     bodyText, err := ioutil.ReadAll(resp.Body)
//     s := string(bodyText)
//     return s
// }
