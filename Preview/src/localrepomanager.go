package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
)

type newBranchHandler interface {
	NewBranch(branch string, dir string)
}

type localRepoManager struct {
	currentBranch    string
	repoSourceDir    string
	localRootDir     string
	repo             *gitlayer
	newBranchObs     newBranchHandler
	enableBranchMode bool
}

func createLocalRepoManager(rootDir string, newBranch newBranchHandler, enableBranchMode bool) *localRepoManager {
	var lrm = &localRepoManager{currentBranch: "master", localRootDir: rootDir}
	lrm.newBranchObs = newBranch
	lrm.enableBranchMode = enableBranchMode

	os.RemoveAll(rootDir) // ignore error since it may not exist
	lrm.repoSourceDir = lrm.ensureDir("source")
	return lrm
}

func (lrm *localRepoManager) ensureDir(subDir string) string {
	var currentPath = path.Join(lrm.localRootDir, subDir)
	var _, err = os.Stat(currentPath)
	if err != nil {
		err = os.MkdirAll(currentPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Couldn't create sourceDir: %v", err.Error())
		}
	}

	return currentPath
}

func (lrm *localRepoManager) getSourceDir() string {
	return lrm.repoSourceDir
}

func (lrm *localRepoManager) getRenderDir() string {
	if lrm.enableBranchMode {
		branchName := lrm.legalizeBranchName(lrm.currentBranch)
		return lrm.ensureDir(branchName)
	}
	return lrm.ensureDir("output")
}

func (lrm *localRepoManager) legalizeBranchName(name string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(name, "")
}

func (lrm *localRepoManager) initialClone(repo string, repopat string) error {
	//TODO: this function should ensure branch name is correct
	fmt.Printf("Initial clone for\n repo: %v\n local dir:%v", repo, lrm.repoSourceDir)
	if repopat != "" {
		fmt.Printf(" with Personal Access Token.\n")
	} else {
		fmt.Printf(" with no authentication.\n")
	}

	re, err := clone(repo, lrm.repoSourceDir, repopat)
	if err != nil {
		fmt.Printf("Error in initial clone: %v\n", err.Error())
		os.Exit(1)
	}
	lrm.repo = re
	fmt.Println("Clone Done.")
	return err
}

func (lrm *localRepoManager) switchBranch(branch string) error {
	if branch != lrm.currentBranch {
		fmt.Printf("Fetching\n")

		err := lrm.repo.checkout(branch)
		if err != nil {
			return fmt.Errorf("checkout failed: %v", err.Error())
		}

		lrm.currentBranch = branch
	}

	err := lrm.repo.pull(branch)
	if err != nil {
		return fmt.Errorf("pull failed: %v", err.Error())
	}
	return nil
}

func (lrm *localRepoManager) handleWebhook(branch string, runjek bool, sendNotify bool) {
	lrm.switchBranch(branch)

	renderDir := lrm.getRenderDir()
	// todo handle branch change
	//TODO run jekyll

	if lrm.enableBranchMode && sendNotify && lrm.newBranchObs != nil {
		lrm.newBranchObs.NewBranch(lrm.legalizeBranchName(branch), renderDir)
	}
}
