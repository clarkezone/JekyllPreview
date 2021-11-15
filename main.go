package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/clarkezone/hookserve/hookserve"

	"github.com/clarkezone/go-execobservable"
)

const (
	reponame          = "JEKPREV_REPO"
	repopat           = "JEKPREV_REPO_PAT"
	webhooksecretname = "JEKPREV_WH_SECRET"
	localdirname      = "JEKPREV_LOCALDIR"
	monitorcmdname    = "JEKPREV_monitorCmd"
	sslwhitelistname  = "JEKPREV_sslwhitelist"
)

var (
	lrm              *localRepoManager
	enableBranchMode bool
	serve            bool
	runjekyll        bool
	sharemgn         *httpShareManager
	authman          *AuthManager
)

type cleanupfunc func()

func main() {
	enableBranchMode = true

	// Read and verify flags
	flag.BoolVar(&serve, "serve", true, "start fileserver")
	flag.BoolVar(&runjekyll, "jekyll", true, "call jekyll")
	flag.Parse()

	repo, repopat, localRootDir, secret, _, sslwhitelist := readEnv()

	verifyFlags(repo, secret, localRootDir)

	// if repopat == "" {
	// 	//delete this
	// }

	authman = NewAuthManager()
	sharemgn = createShareManager(authman.authmux, sslwhitelist)

	// Create Local Repo manager
	lrm = createLocalRepoManager(localRootDir, sharemgn, enableBranchMode)

	//cleanupDone := handleSig(func() { os.RemoveAll(localRootDir) })
	//_ = handleSig(func() { os.RemoveAll(localRootDir) })

	err := lrm.initialClone(repo, repopat)

	initializeJekyll(err)

	startWebhookListener(secret)

	if serve {
		if enableBranchMode {
			//sharemgn.shareBranchPath(lrm.getCurrentBranch(), lrm.getRenderDir())
			sharemgn.shareBranchSubdomain(lrm.getCurrentBranch(), lrm.getRenderDir())
		} else {
			sharemgn.shareRootDir(lrm.getRenderDir())
		}
		//sharemgn.start()
		sharemgn.startsubdomain()
	}

	ch := make(chan bool)
	<-ch
	//<-cleanupDone
}

func verifyFlags(repo string, secret string, localRootDir string) {
	if repo == "" {
		fmt.Printf("Repo must be provided in %v\n", reponame)
		os.Exit(1)
	}

	if secret == "" {
		fmt.Printf("Secret must be provided in %v\n", webhooksecretname)
		os.Exit(1)
	}

	if localRootDir == "" {
		fmt.Printf("Localdir be provided in %v\n", localdirname)
		os.Exit(1)
	}
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
	sslwhitelist := os.Getenv(sslwhitelistname)
	return repo, repopat, localdr, secret, monitorcmdline, sslwhitelist
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
