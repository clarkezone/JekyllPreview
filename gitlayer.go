package main

import (
	"fmt"

	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type gitlayer struct {
	repo *git.Repository
	wt   *git.Worktree
}

func clone(repo string, localfolder string) (*gitlayer, error) {
	gl := &gitlayer{}
	os.RemoveAll(localfolder) // ignore error since it may not exist

	err := os.Mkdir(localfolder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	re, err := git.PlainClone(localfolder, false, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})
	gl.repo = re

	wt, err := gl.repo.Worktree()
	if err != nil {
		fmt.Printf("Get worktree %v\n", err.Error())
		return nil, err
	}
	gl.wt = wt

	if err != nil {
		return nil, err
	}

	return gl, nil
}

func (gl *gitlayer) checkout(branch string) error {

	remote, err := gl.repo.Remote("origin")
	if err != nil {
		fmt.Printf("Get remote %v\n", err.Error())
		return err
	}

	err = remote.Fetch(&git.FetchOptions{})
	if err != nil && err.Error() != "already up-to-date" {
		fmt.Printf("Fetch failed %v\n", err.Error())
		return err
	}

	nm := plumbing.NewBranchReferenceName(branch)

	fmt.Printf("Checking out new branch %v\n", nm)
	err = gl.wt.Checkout(&git.CheckoutOptions{Branch: nm})

	if err != nil {
		fmt.Printf("Checkout new branch failed %v\n", err.Error())
		return err
	}
	return nil
}

func (gl *gitlayer) pull() error {
	err := gl.wt.Pull(&git.PullOptions{})

	if err != nil {
		fmt.Printf("Pull failed %v\n", err.Error())
		return err
	}
	return nil
}
