package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	testreponame         = "TEST_JEKPREV_REPO_NOAUTH"
	testlocaldirname     = "TEST_JEKPREV_LOCALDIR"
	testbranchswitchname = "TEST_JEKPREV_BRANCHSWITCH"
)

func TestCloneNoAuth(t *testing.T) {
	reponame, dirname, _ := getenv()

	os.RemoveAll(dirname)

	clone(reponame, dirname)

	if _, err := os.Stat(dirname); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Clone failed %v\n", err.Error())
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

func TestPullBranch(t *testing.T) {
	reponame, dirname, branch := getenv()

	os.RemoveAll(dirname)

	repo, err := clone(reponame, dirname)
	if err != nil {
		log.Fatal("clone failed")
	}

	err = repo.checkout(branch)
	if err != nil {
		log.Fatal("checkout failed")
	}

	err = repo.pull()
	if err != nil {
		log.Fatal("pull failed")
	}

	os.RemoveAll(dirname)
}

func TestAuth(t *testing.T) {
	//repo, localdr := readEnvTest()
	//clone(repo, localdr)
}

func TestReadEnvTest(t *testing.T) {
	repo, localdr, testbranchswitch := getenv()
	if repo == "" || localdr == "" || testbranchswitch == "" {
		log.Fatalf("Test environment variables not configured: %v\n", testreponame)
	}
}

func getenv() (string, string, string) {
	repo := os.Getenv(testreponame)
	localdr := os.Getenv(testlocaldirname)
	testbranchswitch := os.Getenv(testbranchswitchname)
	return repo, localdr, testbranchswitch
}
