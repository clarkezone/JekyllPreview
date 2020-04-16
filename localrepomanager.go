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

type LocalRepoManager struct {
	currentBranch string
	repoSourceDir string
	localRootDir  string
	repo          *gitlayer
	newBranchObs  newBranchHandler
}

func CreateLocalRepoManager(rootDir string, newBranch newBranchHandler) *LocalRepoManager {
	var lrm = &LocalRepoManager{currentBranch: "master", localRootDir: rootDir}
	lrm.newBranchObs = newBranch

	os.RemoveAll(rootDir) // ignore error since it may not exist
	lrm.repoSourceDir = lrm.ensureDir("source")
	return lrm
}

func (lrm *LocalRepoManager) ensureDir(subDir string) string {
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

func (lrm *LocalRepoManager) getSourceDir() string {
	return lrm.repoSourceDir
}

func (lrm *LocalRepoManager) getCurrentBranch() string {
	return lrm.currentBranch
}

func (lrm *LocalRepoManager) getCurrentBranchRenderDir() string {
	branchName := lrm.legalizeBranchName(lrm.currentBranch)
	return lrm.ensureDir(branchName)
}

func (lrm *LocalRepoManager) legalizeBranchName(name string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(name, "")
}

func (lrm *LocalRepoManager) initialClone(repo string, repopat string) error {
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

func (lrm *LocalRepoManager) handleWebhook(branch string, runjek bool, sendNotify bool) {
	if branch != lrm.currentBranch {
		fmt.Printf("Fetching\n")

		err := lrm.repo.checkout(branch)
		if err != nil {
			log.Fatalf("checkout failed: %v", err.Error())
		}

		lrm.currentBranch = branch
	}

	err := lrm.repo.pull(branch)
	if err != nil {
		log.Fatalf("pull failed: %v", err.Error())
	}

	renderDir := lrm.getCurrentBranchRenderDir()

	if runjek {
		jekBuild(lrm.repoSourceDir, renderDir)

	}

	if sendNotify && lrm.newBranchObs != nil {
		lrm.newBranchObs.NewBranch(lrm.legalizeBranchName(branch), renderDir)
	}
}
