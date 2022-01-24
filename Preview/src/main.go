package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"

	"github.com/clarkezone/hookserve/hookserve"
)

const (
	reponame          = "JEKPREV_REPO"
	repopat           = "JEKPREV_REPO_PAT"
	webhooksecretname = "JEKPREV_WH_SECRET"
	localdirname      = "JEKPREV_LOCALDIR"
	monitorcmdname    = "JEKPREV_monitorCmd"
	initialbranchname = "JEKPREV_initialBranchName"
)

var (
	lrm              *localRepoManager
	enableBranchMode bool
)

type cleanupfunc func()

var serve bool
var runjekyll bool
var webhooklisten bool
var initialclone bool
var sharemgn *httpShareManager

func main() {
	enableBranchMode = false

	// Read and verify flags
	flag.BoolVar(&serve, "serve", false, "start fileserver")
	flag.BoolVar(&runjekyll, "jekyll", false, "call jekyll")
	flag.BoolVar(&webhooklisten, "webhooklisten", false, "listen for webhook messages")
	flag.BoolVar(&initialclone, "initialclone", false, "clone repo")
	flag.Parse()

	//repo, repopat, localRootDir, secret, _ := readEnv()
	repo, _, localRootDir, _, _, initalBranchName := readEnv()

	log.Printf("Called with\nrepo:%v\nlocalRootDir:%v\ninitialclone:%v\nwebhooklisten:%v\nrunjekyll:%v\nserve:%v\n",
		repo, localRootDir,
		initialclone, webhooklisten, runjekyll, serve)

	err := PerformActions(repo, localRootDir, initalBranchName)
	if err != nil {
		log.Printf("Error: %v", err)
		//os.Exit(1)
	}

	ch := make(chan bool)
	<-ch
	//<-cleanupDone
}

func PerformActions(repo string, localRootDir string, initialBranch string) error {
	if serve || runjekyll || webhooklisten || initialclone {
		result := verifyFlags(repo, localRootDir)
		if result != nil {
			return result
		}
	} else {
		return nil
	}

	sourceDir := path.Join(localRootDir, "sourceroot")
	fileinfo, res := os.Stat(sourceDir)
	if fileinfo != nil && res == nil {
		err := os.RemoveAll(sourceDir)
		if err != nil {
			return err
		}
	}

	lrm = createLocalRepoManager(localRootDir, sharemgn, enableBranchMode)

	if initialclone {
		err := lrm.initialClone(repo, repopat)
		if err != nil {
			return err
		}

		if initialBranch != "" {
			return lrm.switchBranch(initialBranch)
		}

	}
	return nil
}

func verifyFlags(repo string, localRootDir string) error {
	if repo == "" {
		return errors.New(fmt.Sprintf("Repo must be provided in %v\n", reponame))
	}

	if localRootDir == "" {
		return errors.New(fmt.Sprintf("Localdir be provided in %v\n", localRootDir))
	} else {
		fileinfo, res := os.Stat(localRootDir)
		if res != nil {
			return errors.New(fmt.Sprintf("Localdir must exist %v\n", localRootDir))
		}
		if !fileinfo.IsDir() {
			return errors.New(fmt.Sprintf("Localdir must be a directory %v\n", localRootDir))
		}
	}
	return nil
}

func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
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

func readEnv() (string, string, string, string, string, string) {
	repo := os.Getenv(reponame)
	repopat := os.Getenv(repopat)
	localdr := os.Getenv(localdirname)
	secret := os.Getenv(webhooksecretname)
	monitorcmdline := os.Getenv(monitorcmdname)
	initalbranchname := os.Getenv(initialbranchname)
	return repo, repopat, localdr, secret, monitorcmdline, initalbranchname
}

func startWebhookListener(secret string) {
	go func() {
		fmt.Printf("Monitoring started\n")
		server := hookserve.NewServer()
		server.Port = 8080
		server.Secret = secret
		server.GoListenAndServe()

		for event := range server.Events {
			fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
			lrm.handleWebhook(event.Branch, runjekyll, runjekyll)
		}
	}()
}
