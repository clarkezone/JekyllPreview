package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/clarkezone/hookserve/hookserve"

	"github.com/clarkezone/go-execobservable"
)

const (
	reponame          = "JEKPREV_REPO"
	repopat           = "JEKPREV_REPO_PAT"
	webhooksecretname = "JEKPREV_WH_SECRET"
	localdirname      = "JEKPREV_LOCALDIR"
	monitorcmdname    = "JEKPREV_monitorCmd"
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
	repo, _, localRootDir, _, _ := readEnv()

	log.Printf("Called with\nrepo:%v\nlocalRootDir:%v\ninitialclone:%v\nwebhooklisten:%v\nrunjekyll:%v\nserve:%v\n",
		repo, localRootDir,
		initialclone, webhooklisten, runjekyll, serve)

	err := PerformActions(repo, localRootDir)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	//	sharemgn = createShareManager()
	//
	//
	//	//cleanupDone := handleSig(func() { os.RemoveAll(localRootDir) })
	//	//_ = handleSig(func() { os.RemoveAll(localRootDir) })
	//
	//
	//	initializeJekyll(err)
	//
	//	startWebhookListener(secret)
	//
	//	if serve {
	//		if enableBranchMode {
	//			sharemgn.shareBranch(lrm.getCurrentBranch(), lrm.getRenderDir())
	//		} else {
	//			sharemgn.shareRootDir(lrm.getRenderDir())
	//		}
	//		sharemgn.start()
	//	}
	//
	//	ch := make(chan bool)
	//	<-ch
	//	//<-cleanupDone
}

func PerformActions(repo string, localRootDir string) error {
	if serve || runjekyll || webhooklisten || initialclone {
		result := verifyFlags(repo, localRootDir)
		if result != nil {
			return result
		}
	} else {
		return nil
	}

	lrm = createLocalRepoManager(localRootDir, sharemgn, enableBranchMode)

	if initialclone {
		return lrm.initialClone(repo, repopat)
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
		empty, err := IsEmpty(localRootDir)
		if !empty {
			return errors.New(fmt.Sprintf("Localdir must be empty %v\n", localRootDir))
		}

		if err != nil {
			return err
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

func readEnv() (string, string, string, string, string) {
	repo := os.Getenv(reponame)
	repopat := os.Getenv(repopat)
	localdr := os.Getenv(localdirname)
	secret := os.Getenv(webhooksecretname)
	monitorcmdline := os.Getenv(monitorcmdname)
	return repo, repopat, localdr, secret, monitorcmdline
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

func initializeJekyll(err error) {
	if runjekyll {
		fmt.Printf("Starting Jekyll with sourcedir %v..\n", lrm.getSourceDir())
		err = jekPrepare(lrm.getSourceDir())
		if err != nil {
			fmt.Printf("Error in Jekyll prep: %v\n", err.Error())
			os.Exit(1)
		}

		// Note jekyll build errors are truncated by exec so you only see the warning line
		// not the actual error.  Use the streaming cmdversion to show complete spew
		err = jekBuild(lrm.getSourceDir(), lrm.getRenderDir())
		if err != nil {
			fmt.Printf("Error in Jekyll build: %v\n", err.Error())
			os.Exit(1)
		}
	}
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

	cmdstring := "bundle install"

	output := &outputprogress{}

	runner := &execobservable.CmdRunner{}
	runner.RunCommand("sh", localfolder, output, "-c", cmdstring)

	return nil
}

func jekBuild(localfolder string, outputfolder string) error {
	cmd := exec.Command("chown", "-R", "jekyll:jekyll", outputfolder)
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Unable to change ownership")
	}

	//cmd := exec.Command("bundle exec jekyll build --destination " + outputfolder)
	fmt.Printf("Running jekyll with sourcedir %v and output %v\n", localfolder, outputfolder)
	// cmd := exec.Command("bundle", "exec", "jekyll", "build", "--destination", outputfolder)
	// var errString bytes.Buffer
	// cmd.Stderr = &errString
	// cmd.Dir = localfolder
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Printf("Error: %q\n", errString.String())
	// 	return err
	// }

	cmdstring := "bundle exec jekyll build --destination " + outputfolder

	output := &outputprogress{}

	runner := &execobservable.CmdRunner{}
	runner.RunCommand("sh", localfolder, output, "-c", cmdstring)

	return nil
}

type outputprogress struct {
}

func (t outputprogress) Progress(s string, sr execobservable.SendResponse) {
	fmt.Printf("%v", s)
}
