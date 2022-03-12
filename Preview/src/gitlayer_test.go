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
	testsecurereponame   = "TEST_JEKPREV_SECURE_REPO_NOAUTH"
	testsecureclonepw    = "TEST_JEKPREV_SECURECLONEPW"
)

//configure environment variables by:
// 1. command palette: open settings (json)
// 2. append the following
//"go.testEnvVars": {
//	"TEST_JEKPREV_REPO_NOAUTH": "https://URL:
//	"TEST_JEKPREV_LOCALDIR": "/tmp/jekpreview_test",
//	"TEST_JEKPREV_BRANCHSWITCH": "testbranch",
//	"TEST_JEKPREV_SECURE_REPO_NOAUTH": "https://",
//	"TEST_JEKPREV_SECURECLONEPW": "unused",
//  },

func TestAllReadEnvTest(t *testing.T) {
	t.Logf("TestAllReadEnvTest")
	repo, localdr, testbranchswitch, securereponame, secureclonepw := getenv()
	if repo == "" || localdr == "" || testbranchswitch == "" || securereponame == "" || secureclonepw == "" {
		log.Fatalf("Test environment variables not configured")
	}
}

func TestCloneNoAuth(t *testing.T) {
	t.Logf("TestCloneNoAuth")
	reponame, dirname, _, _, _ := getenv()

	os.RemoveAll(dirname)

	clone(reponame, dirname, "")

	if _, err := os.Stat(dirname); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Clone failed %v\n", err.Error())
		}
	}

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatalf("TestCloneNoAuth: clone failed %v", err.Error())
	}

	if len(infos) < 8 {
		log.Fatalf("TestCloneNoAuth: clone failed expected %v, found %v", 9, len(infos))
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
	t.Logf("TestPullBranch")
	reponame, dirname, branch, _, _ := getenv()

	os.RemoveAll(dirname)

	repo, err := clone(reponame, dirname, "")
	if err != nil {
		log.Fatal("TestPullBranch: clone failed")
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

	if len(infos) != 12 { //One extra for .git
		log.Fatalf("pull failed file mismatch error expected 9 found %v", len(infos))
	}

	os.RemoveAll(dirname)
}

func TestCloneAuth(t *testing.T) {
	t.Logf("TestCloneAuth")
	_, dirname, _, secureproname, pw := getenv()
	//reponame, dirname, branch, pw := getenv()

	if pw == "unused" {
		return
	}

	os.RemoveAll(dirname)

	_, err := clone(secureproname, dirname, pw)
	//repo, err := clone(reponame, dirname, "", pw)
	if err != nil {
		log.Fatal("TestCloneAuth: clone failed")
	}

	// err = repo.checkout(branch)
	// if err != nil {
	// 	log.Fatal("checkout failed")
	// }

	// err = repo.pull(branch)
	// if err != nil {
	// 	log.Fatal("pull failed")
	// }

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal("pull failed")
	}

	if len(infos) != 3 { //One extra for .git
		log.Fatalf("pull failed file mismatch error")
	}

	os.RemoveAll(dirname)
}

func getenv() (string, string, string, string, string) {
	repo := os.Getenv(testreponame)
	localdr := os.Getenv(testlocaldirname)
	testbranchswitch := os.Getenv(testbranchswitchname)
	reposecure := os.Getenv(testsecurereponame)
	secureclonepw := os.Getenv(testsecureclonepw)
	return repo, localdr, testbranchswitch, reposecure, secureclonepw
}
