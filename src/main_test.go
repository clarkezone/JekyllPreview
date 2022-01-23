package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestVerifyInitialClone(t *testing.T) {
	repo, localdr, _, _, _ := getenv()
	initialclone = true
	localdr, err := ioutil.TempDir("/tmp", "jekylltest")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(localdr)
	err = PerformActions(repo, localdr)

	if !containsItems(path.Join(localdr, "source")) {
		t.Error("no items cloned")
	}

	if err != nil {
		fmt.Printf(err.Error())
		t.Fail()
	}
}

func containsItems(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	names, err := f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return false
	}
	if len(names) < 1 {
		return false
	}
	return true
}
