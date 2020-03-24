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

func TestAllReadEnvTest(t *testing.T) {
	repo, localdr, testbranchswitch := getenv()
	if repo == "" || localdr == "" || testbranchswitch == "" {
		log.Fatalf("Test environment variables not configured: %v\n", testreponame)
	}
}

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

// func TestPullSameBranch(t *testing.T) {
// 	//_, dirname, branch := getenv()
// 	reponame, dirname, _ := getenv()
// 	os.RemoveAll(dirname)

// 	repo, err := clone(reponame, dirname)
// 	//repo, err := open(dirname)
// 	if err != nil {
// 		log.Fatal("open failed")
// 	}

// 	err = repo.checkout("debugsinglepull")

// 	//err = repo.pull(branch)
// 	if err != nil {
// 		log.Fatal("pull failed")
// 	}
// }

// func TestPullSameBranchPull(t *testing.T) {
// 	_, dirname, _ := getenv()
// 	//reponame, dirname, _ := getenv()
// 	//os.RemoveAll(dirname)

// 	//repo, err := clone(reponame, dirname)
// 	repo, err := open(dirname)
// 	if err != nil {
// 		log.Fatal("open failed")
// 	}

// 	//err = repo.checkout("debugsinglepull")

// 	err = repo.pull("debugsinglepull")
// 	if err != nil {
// 		log.Fatal("pull failed")
// 	}
// }

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

	err = repo.pull(branch)
	if err != nil {
		log.Fatal("pull failed")
	}

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal("pull failed")
	}

	if len(infos) != 3 { //One extra for .git
		log.Fatalf("pull failed file mismatch error")
	}

	os.RemoveAll(dirname)
}

func TestAuth(t *testing.T) {
	//repo, localdr := readEnvTest()
	//clone(repo, localdr)
}

func getenv() (string, string, string) {
	repo := os.Getenv(testreponame)
	localdr := os.Getenv(testlocaldirname)
	testbranchswitch := os.Getenv(testbranchswitchname)
	return repo, localdr, testbranchswitch
}
