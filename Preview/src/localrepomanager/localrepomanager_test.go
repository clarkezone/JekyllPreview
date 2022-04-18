package localrepomanager

import (
	"io/ioutil"
	"os"
	"testing"
)

var lrm *LocalRepoManager

func SkipCI(t *testing.T) {
	if os.Getenv("TEST_JEKPREV_TESTLOCALK8S") == "" {
		t.Skip("Skipping K8slocaltest")
	}
}

func TestSourceDir(t *testing.T) {
	lrm = CreateLocalRepoManager("test", nil, true, nil)

	res := lrm.getSourceDir()

	if res != "test/source" {
		t.Fatalf("Incorrect source dir")
	}

	os.RemoveAll("test")
}

func TestCreateLocalRepoManager(t *testing.T) {
	_ = CreateLocalRepoManager("test", nil, true, nil)

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
	lrm := CreateLocalRepoManager("test", nil, true, nil)
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
	lrm := CreateLocalRepoManager("test", nil, true, nil)

	dir := lrm.getRenderDir()

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
	SkipCI(t)
	repo, dirname, _, _, _ := Getenv()

	lrm := CreateLocalRepoManager(dirname, nil, true, nil)
	err := lrm.InitialClone(repo, "")
	if err != nil {
		t.Fatalf("error in initial clonse")
	}

	os.RemoveAll(dirname)
}

// TODO unregister

// func TestLRMSwitchBranch(t *testing.T) {
// 	_, dirname, branch, secureRepo, pat := getenv()

// 	lrm := CreateLocalRepoManager(dirname)
// 	lrm.initialClone(secureRepo, pat)

// 	lrm.handleWebhook(branch, false, true)

// 	branchDir := lrm.getCurrentBranchRenderDir()

// 	if branchDir != path.Join(dirname, branch) {
// 		t.Fatalf("incorrect new dir")
// 	}

// 	os.RemoveAll(dirname)
// }

//func TestLRMSwitchBranchBackToMain(t *testing.T) {
//	_, dirname, branch, secureRepo, pat := getenv()
//
//	sharemgn := createShareManager()
//
//	lrm := CreateLocalRepoManager(dirname, sharemgn, true)
//	lrm.initialClone(secureRepo, pat)
//
//	lrm.handleWebhook(branch, false, true)
//
//	branchDir := lrm.getRenderDir()
//
//	if branchDir != path.Join(dirname, branch) {
//		t.Fatalf("incorrect new dir")
//	}
//
//	lrm.handleWebhook("master", false, true)
//
//	os.RemoveAll(dirname)
//}
