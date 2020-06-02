package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func _TestInitShareManager(t *testing.T) {
	sm := createShareManager(nil)

	sm.start()
	sm.shareBranchPath("/master/", "./test/one")
	sm.shareBranchPath("/pepper/", "./test/two")

	c := http.Client{}
	resp, err := c.Get("http://localhost:8085/master")
	if err != nil {
		t.Fatalf("Request failed")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Request failed")
	}

	result := string(bytes)
	if !strings.HasPrefix(result, "<headone>") {
		t.Fatalf("First doesn't match %v", result)
	}

}

func _TestInitShareManagerSubdomain(t *testing.T) {
	sm := createShareManager(nil)

	sm.startsubdomain()
	sm.shareBranchSubdomain("master", "./test/one")
	sm.shareBranchSubdomain("pepper", "./test/two")

	c := http.Client{}
	resp, err := c.Get("http://master.localhost:8085/")
	if err != nil {
		t.Fatalf("Request failed")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Request failed")
	}

	result := string(bytes)
	if !strings.HasPrefix(result, "<headone>") {
		t.Fatalf("First doesn't match %v", result)
	}

	resp, err = c.Get("http://pepper.localhost:8085/")
	if err != nil {
		t.Fatalf("Request failed")
	}

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Request failed")
	}

	result = string(bytes)
	if !strings.HasPrefix(result, "<headtwo>") {
		t.Fatalf("Second doesn't match")
	}

}
