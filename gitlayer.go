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
	pat  string
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
	if err != nil {
		fmt.Printf("Plainclone %v\n", err.Error())
		return nil, err
	}
	gl.repo = re

	wt, err := gl.repo.Worktree()
	if err != nil {
		fmt.Printf("Get worktree %v\n", err.Error())
		return nil, err
	}
	gl.wt = wt
	gl.pat = pw

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

	var feo *git.FetchOptions

	if gl.pat == "" {
		feo = &git.FetchOptions{}
	} else {
		feo = &git.FetchOptions{
			Auth: &http.BasicAuth{
				Username: "abc123", // yes, this can be anything except an empty string
				Password: gl.pat,
			},
		}
	}

	err = remote.Fetch(feo)
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

	var feo *git.PullOptions

	if gl.pat == "" {
		feo = &git.PullOptions{ReferenceName: nm}
	} else {
		feo = &git.PullOptions{
			ReferenceName: nm,
			Auth: &http.BasicAuth{
				Username: "abc123", // yes, this can be anything except an empty string
				Password: gl.pat,
			},
		}
	}

	err := gl.wt.Pull(feo)

	if err != nil && err.Error() != "already up-to-date" {
		fmt.Printf("Pull failed %v\n", err.Error())
		return err
	}
	return nil
}

func (gl *gitlayer) getBranch() (string, error) {
	//Note this approach doesn't work
	return GetCurrentBranchFromRepository(gl.repo)
}

func GetCurrentBranchFromRepository(repository *git.Repository) (string, error) {
	branchRefs, err := repository.Branches()
	if err != nil {
		return "", err
	}

	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}

	var currentBranchName string
	err = branchRefs.ForEach(func(branchRef *plumbing.Reference) error {
		if branchRef.Hash() == headRef.Hash() {
			currentBranchName = branchRef.Name().Short()

			return nil
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return currentBranchName, nil
}
