package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"

	"github.com/clarkezone/go-execobservable"
	"github.com/phayes/hookserve/hookserve"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	reponame       = "JEKPREV_REPO"
	secretname     = "JEKPREV_SECRET"
	localdirname   = "JEKPREV_LOCALDIR"
	monitorcmdname = "JEKPREV_monitorCmd"
)

type cleanupfunc func()

var serve bool
var runjekyll bool

func main() {
	flag.BoolVar(&serve, "serve", true, "start fileserver")
	flag.BoolVar(&runjekyll, "jekyll", true, "call jekyll")
	flag.Parse()

	repo, localdir, secret, _ := readEnv()

	if repo == "" {
		fmt.Printf("Repo must be provided in %v\n", reponame)
		os.Exit(1)
	}

	if secret == "" {
		fmt.Printf("Secret must be provided in %v\n", secretname)
		os.Exit(1)
	}

	if localdir == "" {
		fmt.Printf("Localdir be provided in %v\n", localdirname)
		os.Exit(1)
	}

	//cleanupDone := handleSig(func() { os.RemoveAll(localdir) })
	//_ = handleSig(func() { os.RemoveAll(localdir) })

	fmt.Printf("Initial clone for\n repo: %v\n local dir:%v\n", repo, localdir)

	re, err := clone(repo, localdir)
	if err != nil {
		fmt.Printf("Error in initial clone: %v\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Clone Done.")

	if runjekyll {
		err = jekPrepare(localdir)
		if err != nil {
			fmt.Printf("Error in Jekyll prep: %v\n", err.Error())
			os.Exit(1)
		}

		err = jekBuild(localdir, "/srv/jekyll/output/master")
		if err != nil {
			fmt.Printf("Error in Jekyll build: %v\n", err.Error())
			os.Exit(1)
		}
	}

	go func() {
		fmt.Printf("Monitoring started\n")
		err := monitor(secret, localdir, re)
		if err != nil {
			fmt.Printf("Monitor failed: %v\n", err.Error())
			os.Exit(1)
		}
	}()

	if serve {
		http.Handle("/", http.FileServer(http.Dir("/srv/jekyll/output/master")))
		http.ListenAndServe(":8085", nil)
	}

	//<-cleanupDone

}

func handleSig(cleanupwork cleanupfunc) chan struct{} {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Printf("\nReceived an interrupt, stopping services...\n")
		if cleanupwork != nil {
			cleanupwork()
		}

		close(cleanupDone)
	}()
	return cleanupDone
}

func readEnv() (string, string, string, string) {
	repo := os.Getenv(reponame)
	localdr := os.Getenv(localdirname)
	secret := os.Getenv(secretname)
	monitorcmdline := os.Getenv(monitorcmdname)
	return repo, localdr, secret, monitorcmdline
}

func monitor(secret string, localfolder string, repo *git.Repository) error {
	currentBranch := "master"
	w, err := repo.Worktree()
	//fmt.Printf("Current branch from git %v\n")
	server := hookserve.NewServer()
	server.Port = 8080
	server.Secret = secret
	server.GoListenAndServe()

	// Everytime the server receives a webhook event, print the results
	for event := range server.Events {
		fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)

		if event.Branch != currentBranch {
			fmt.Printf("Fetching\n")

			remote, err := repo.Remote("origin")
			if err != nil {
				fmt.Printf("Get remote %v\n", err.Error())
				return err
			}

			err = remote.Fetch(&git.FetchOptions{})
			if err != nil {
				fmt.Printf("Fetch failed %v\n", err.Error())
				return err
			}

			nm := plumbing.NewBranchReferenceName(event.Branch)
			fmt.Printf("Checking out new branch %v\n", nm)
			err = w.Checkout(&git.CheckoutOptions{Branch: nm})

			if err != nil {
				fmt.Printf("Checkout new branch failed %v\n", err.Error())
				return err
			}

			currentBranch = event.Branch

			jekBuild(localfolder, "/srv/jekyll/output/master")
		}

		fmt.Printf("Pull branch: %v\n", event.Branch)
		err = w.Pull(&git.PullOptions{})

		if err != nil {
			fmt.Printf("Pull failed %v\n", err.Error())
			return err
		}

	}
	return nil
}

func clone(repo string, localfolder string) (*git.Repository, error) {
	os.RemoveAll(localfolder) // ignore error since it may not exist

	err := os.Mkdir(localfolder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	re, err := git.PlainClone(localfolder, false, &git.CloneOptions{
		URL:      repo,
		Progress: os.Stdout,
	})

	if err != nil {
		return re, err
	}

	return re, nil
}

func jekPrepare(localfolder string) error {
	cmd := exec.Command("bundle", "install")
	var errString bytes.Buffer
	cmd.Stderr = &errString
	cmd.Dir = localfolder
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %q\n", errString.String())
		return err
	}
	return nil
}

func jekBuild(localfolder string, outputfolder string) error {
	//cmd := exec.Command("bundle exec jekyll build --destination " + outputfolder)
	cmd := exec.Command("bundle", "exec", "jekyll", "build", "--destination", outputfolder)
	var errString bytes.Buffer
	cmd.Stderr = &errString
	cmd.Dir = localfolder
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %q\n", errString.String())
		return err
	}
	return nil
}

type outputprogress struct {
}

func (t outputprogress) Progress(s string, sr execobservable.SendResponse) {
	fmt.Printf("%v", s)
}
