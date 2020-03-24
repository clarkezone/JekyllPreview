package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5/plumbing"
)

type gitlayer struct {
	repo *git.Repository
	wt   *git.Worktree
}

func open(localfolder string) (*gitlayer, error) {
	gl := &gitlayer{}
	re, err := git.PlainOpen(localfolder)
	if err != nil {
		return nil, err
	}
	wt, err := re.Worktree()
	gl.repo = re
	gl.wt = wt
	return gl, nil
}

func clone(repo string, localfolder string, pw string) (*gitlayer, error) {
	gl := &gitlayer{}
	os.RemoveAll(localfolder) // ignore error since it may not exist

	err := os.Mkdir(localfolder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var clo *git.CloneOptions

	if pw == "" {
		clo = &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
		}
	} else {
		clo = &git.CloneOptions{
			URL:      repo,
			Progress: os.Stdout,
			Auth: &http.BasicAuth{
				Username: "abc123", // yes, this can be anything except an empty string
				Password: pw,
			},
		}
	}

	re, err := git.PlainClone(localfolder, false, clo)
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

	nm := plumbing.NewRemoteReferenceName(remote.Config().Name, branch)

	fmt.Printf("Checking out new branch %v\n", nm)
	err = gl.wt.Checkout(&git.CheckoutOptions{Branch: nm})

	if err != nil {
		fmt.Printf("Checkout new branch failed %v\n", err.Error())
		return err
	}
	return nil
}

func (gl *gitlayer) pull(branch string) error {
	fmt.Printf("Pulling branch %v\n", branch)
	nm := plumbing.NewBranchReferenceName(branch)

	err := gl.wt.Pull(&git.PullOptions{ReferenceName: nm})

	if err != nil && err.Error() != "already up-to-date" {
		fmt.Printf("Pull failed %v\n", err.Error())
		return err
	}
	return nil
}
