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
	"github.com/clarkezone/hookserve/hookserve"
)

const (
	reponame          = "JEKPREV_REPO"
	repopat           = "JEKPREV_REPO_PAT"
	webhooksecretname = "JEKPREV_WH_SECRET"
	localdirname      = "JEKPREV_LOCALDIR"
	monitorcmdname    = "JEKPREV_monitorCmd"
)

type cleanupfunc func()

var serve bool
var runjekyll bool

func main() {
	flag.BoolVar(&serve, "serve", true, "start fileserver")
	flag.BoolVar(&runjekyll, "jekyll", true, "call jekyll")
	flag.Parse()

	repo, repopat, localdir, secret, _ := readEnv()

	if repo == "" {
		fmt.Printf("Repo must be provided in %v\n", reponame)
		os.Exit(1)
	}

	if secret == "" {
		fmt.Printf("Secret must be provided in %v\n", webhooksecretname)
		os.Exit(1)
	}

	if localdir == "" {
		fmt.Printf("Localdir be provided in %v\n", localdirname)
		os.Exit(1)
	}

	//cleanupDone := handleSig(func() { os.RemoveAll(localdir) })
	//_ = handleSig(func() { os.RemoveAll(localdir) })

	fmt.Printf("Initial clone for\n repo: %v\n local dir:%v", repo, localdir)
	if repopat != "" {
		fmt.Printf(" with Personal Access Token.\n")
	} else {
		fmt.Printf(" with no authentication.\n")
	}

	re, err := clone(repo, localdir, repopat)
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
		err := monitor(secret, localdir, re, runjekyll)
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

func readEnv() (string, string, string, string, string) {
	repo := os.Getenv(reponame)
	repopat := os.Getenv(repopat)
	localdr := os.Getenv(localdirname)
	secret := os.Getenv(webhooksecretname)
	monitorcmdline := os.Getenv(monitorcmdname)
	return repo, repopat, localdr, secret, monitorcmdline
}

func monitor(secret string, localfolder string, repo *gitlayer, runjek bool) error {
	currentBranch := "master"

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

			repo.checkout(event.Branch)

			currentBranch = event.Branch
		}

		repo.pull(event.Branch)

		if runjek {
			jekBuild(localfolder, "/srv/jekyll/output/master")
		}
	}
	return nil
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
