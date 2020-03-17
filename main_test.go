package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	testreponame     = "TEST_JEKPREV_REPO"
	testlocaldirname = "TEST_JEKPREV_LOCALDIR"
)

func TestCloneNoAuth(t *testing.T) {
	dirname := "test_clarkezone"

	os.RemoveAll(dirname)

	clone("https://github.com/clarkezone/clarkezone.github.io.git", dirname)

	if _, err := os.Stat(dirname); err != nil {
		if os.IsNotExist(err) {
			log.Fatal("Clone failed")
		} else {
			// other error
		}
	}

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal("clone failed")
	}

	if len(infos) < 5 {
		log.Fatalf("clone failed")
	}

	os.RemoveAll(dirname)
}

func TestAuth(t *testing.T) {
	//repo, localdr := readEnvTest()
	//clone(repo, localdr)
}

func readEnvTest() (string, string) {
	repo := os.Getenv(testreponame)
	localdr := os.Getenv(testlocaldirname)
	if repo == "" || localdr == "" {
		log.Fatalf("Test environment variables not configured: %v\n", testreponame)
	}
	return repo, localdr
}
