package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/phayes/hookserve/hookserve"
)

var (
	cloneFlag   = flag.Bool("clone", false, "Perform initial clone")
	monitorFlag = flag.Bool("monitor", false, "Monitor for changes via webhook and pull")
)

func main() {
	flag.Parse()

	repo := "https://github.com/clarkezone/clarkezone.github.io.git"
	branch := "acme-ngrok"
	//localdir := "/srv/jekyll/source"
	localdir := "clarkezone.github.io"
	secret := "onetwothree"

	if *cloneFlag {
		fmt.Printf("Initial clone for\n repo: %v\n branch: %v\n local dir:%v\n", repo, branch, localdir)

		err := clone(repo, branch, localdir)
		if err != nil {
			fmt.Printf("Error in initial clone: %v\n", err.Error())
			os.Exit(1)
		}
		fmt.Println("Done.")
	}

	if *monitorFlag {
		fmt.Printf("Monitoring branch: %v", branch)
		monitor(secret, branch)
	}
}

func monitor(secret string, branch string) {
	server := hookserve.NewServer()
	server.Port = 8080
	server.Secret = secret
	//gitLocalDir: /sourcetmp
	//targetBranch: acme-ngrok
	//gitrepo: gitRepo: https://github.com/clarkezone/clarkezone.github.io.git
	server.GoListenAndServe()

	// Everytime the server receives a webhook event, print the results
	for event := range server.Events {
		fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
		if event.Branch == branch {
			fmt.Println("Match")
		}
	}
}

func clone(repo string, branch string, localfolder string) error {
	os.Mkdir(localfolder, os.ModePerm)

	cmd := exec.Command("git", "clone", repo, ".")
	cmd.Dir = localfolder
	err := cmd.Run()

	if err != nil {
		return err
	}

	//os.Chown("/srv/jekyll/source", 1000, 1000) not recursive

	cmd = exec.Command("git", "checkout", branch)
	cmd.Dir = localfolder
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func prepJekyll(localfolder string, sitdir string) error {
	os.Mkdir(sitdir, os.ModePerm)

	cmd := exec.Command("chown", "-R", "jekyll:jekyll", localfolder)
	err := cmd.Run()

	if err != nil {
		return err
	}

	cmd = exec.Command("chown", "-R", "jekyll:jekyll", sitdir)
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
