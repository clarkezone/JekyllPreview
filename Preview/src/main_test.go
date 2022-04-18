package main

import (
	"io"
	"log"
	"os"
	"path"
	"testing"
)

func TestVerifyInitialCloneDefaultBranch(t *testing.T) {
	t.Logf("TestVerifyInitialCloneDefaultBranch")
	//	repo, _, _, _, _, _ := llrm.ReadEnv()
	repo := "https://github.com/clarkezone/JekyllPreview.git"
	initialclone = true
	localdr, err := t.TempDir()
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(localdr)
	err = PerformActions(repo, localdr, "", false, "testns", false)

	if !containsItems(path.Join(localdr, "source")) {
		t.Error("no items cloned")
	}

	if err != nil {
		t.Fail()
	}
}

func TestVerifyInitialClonewithInitialBranch(t *testing.T) {
	//	repo, _, initialBranch, _, _, _ := llrm.ReadEnv()
	repo := "https://github.com/clarkezone/JekyllPreview.git"
	initialBranch := "BugFix"
	t.Logf("TestVerifyInitialClonewithInitialBranch branch %v", initialBranch)
	initialclone = true
	localdr, err := t.TempDir()
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(localdr)
	err = PerformActions(repo, localdr, initialBranch, false, "testns", false)

	if !containsItems(path.Join(localdr, "source")) {
		t.Error("no items cloned")
	}

	if err != nil {
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
