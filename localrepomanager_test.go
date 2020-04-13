package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSourceDir(t *testing.T) {
	lrm = CreateLocalRepoManager("test")

	res := lrm.getSourceDir()

	if res != "test/source" {
		t.Fatalf("Incorrect source dir")
	}

	os.RemoveAll("test")
}

func TestCreateLocalRepoManager(t *testing.T) {
	_ = CreateLocalRepoManager("test")

	_, err := ioutil.ReadDir("test")
	if err != nil {
		t.Fatalf("Directory didn't get created")
	}

	_, err = ioutil.ReadDir("test/source")
	if err != nil {
		t.Fatalf("Directory didn't get created")
	}

	os.RemoveAll("test")
}

func TestLegalizeBranchName(t *testing.T) {
	lrm := CreateLocalRepoManager("test")
	result := lrm.legalizeBranchName("foo")
	if result != "foo" {
		t.Fatalf("result incorrect")
	}

	result = lrm.legalizeBranchName("f-o-o")
	if result != "foo" {
		t.Fatalf("result incorrect")
	}

	result = lrm.legalizeBranchName("f*o*o")
	if result != "foo" {
		t.Fatalf("result incorrect")
	}

	os.RemoveAll("test")
}

func TestGetCurrentBranchRender(t *testing.T) {
	lrm := CreateLocalRepoManager("test")

	dir := lrm.getCurrentBranchRenderDir()

	if dir != "test/master" {
		t.Fatalf("Wrong name")
	}

	_, err := ioutil.ReadDir("test/master")
	if err != nil {
		t.Fatalf("Directory didn't get created")
	}

	os.RemoveAll("test")
}

func TestLRMCheckout(t *testing.T) {
	_, dirname, _, secureRepo, pat := getenv()

	lrm := CreateLocalRepoManager(dirname)
	err := lrm.initialClone(secureRepo, pat)
	if err != nil {
		t.Fatalf("error in initial clonse")
	}

	os.RemoveAll(dirname)
}
func TestLRMSwitchBranch(t *testing.T) {
	_, dirname, branch, secureRepo, pat := getenv()

	lrm := CreateLocalRepoManager(dirname)
	lrm.initialClone(secureRepo, pat)

	lrm.handleWebhook(branch, false)

	branchDir := lrm.getCurrentBranchRenderDir()

	if branchDir != path.Join(dirname, branch) {
		t.Fatalf("incorrect new dir")
	}

	os.RemoveAll(dirname)
}
