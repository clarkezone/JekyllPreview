package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/phayes/hookserve/hookserve"
)

var (
	cloneFlag   = flag.Bool("clone", false, "Perform initial clone")
	monitorFlag = flag.Bool("monitor", false, "Monitor for changes via webhook and pull")
)

const (
	reponame       = "JEKPREV_REPO"
	secretname     = "JEKPREV_SECRET"
	localdirname   = "JEKPREV_LOCALDIR"
	monitorcmdname = "JEKPREV_monitorCmd"
)

func main() {
	flag.Parse()

	repo, localdir, secret, monitorcmdline := readEnv()

	if repo == "" {
		fmt.Printf("Repo must be provided in %v\n", reponame)
		os.Exit(1)
	}

	if secret == "" {
		fmt.Printf("Repo must be provided in %v\n", secretname)
		os.Exit(1)
	}

	if localdir == "" {
		fmt.Printf("Localdir be provided in %v\n", localdirname)
		os.Exit(1)
	}

	//if *cloneFlag {
	{
		fmt.Printf("Initial clone for\n repo: %v\n local dir:%v\n", repo, localdir)

		err := clone(repo, localdir)
		if err != nil {
			fmt.Printf("Error in initial clone: %v\n", err.Error())
			os.Exit(1)
		}
		fmt.Println("Done.")
	}

	//var comp chan bool

	//if *monitorFlag {
	{
		//comp = make(chan bool)
		go func() {
			fmt.Printf("Monitoring started\n")
			err := monitor(secret, localdir)
			if err != nil {
				fmt.Printf("Monitor failed: %v\n", err.Error())
				os.Exit(1)
			}
		}()
	}

	if monitorcmdline != "" {
		fmt.Printf("Running commandline %v\n", monitorcmdline)
		err := prepJekyll(localdir, "_site")
		if err != nil {
			fmt.Printf("PrepJekyll failed: %v\n", err.Error())
			os.Exit(1)
		}
		err = runJekyllScript(monitorcmdline)
		if err != nil {
			fmt.Printf("Monitor cmdline failed: %v\n", err.Error())
			os.Exit(1)
		}
	}

	//if *monitorFlag {
	// {
	// 	<-comp
	// }
	handleSig(localdir)
}

func cleanup() {

}

func handleSig(localfolder string) {
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Printf("\nReceived an interrupt, stopping services...\n")
		os.RemoveAll(localfolder)
		close(cleanupDone)
	}()
	<-cleanupDone
}

func readEnv() (string, string, string, string) {
	repo := os.Getenv(reponame)
	localdr := os.Getenv(localdirname)
	secret := os.Getenv(secretname)
	monitorcmdline := os.Getenv(monitorcmdname)
	return repo, localdr, secret, monitorcmdline
}

func monitor(secret string, localfolder string) error {
	currentBranch := ""
	server := hookserve.NewServer()
	server.Port = 8080
	server.Secret = secret
	server.GoListenAndServe()

	// Everytime the server receives a webhook event, print the results
	for event := range server.Events {
		fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)

		if event.Branch != currentBranch {
			currentBranch = event.Branch
			fmt.Printf("Checking out new branch %v\n", currentBranch)
			cmd := exec.Command("git", "checkout", currentBranch)
			cmd.Dir = localfolder
			err := cmd.Run()

			if err != nil {
				return err
			}
		}

		fmt.Printf("Pull branch: %v\n", event.Branch)
		cmd := exec.Command("git", "pull")
		cmd.Dir = localfolder
		err := cmd.Run()

		if err != nil {
			return err
		}

	}
	return nil
}

func clone(repo string, localfolder string) error {
	err := os.Mkdir(localfolder, os.ModePerm)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", repo, ".")
	cmd.Dir = localfolder
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func prepJekyll(localfolder string, sitdir string) error {
	cmd := exec.Command("chown", "-R", "jekyll:jekyll", localfolder)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func runJekyllScript(cmdstring string) error {
	cmd := exec.Command("sh", "-c", cmdstring)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
