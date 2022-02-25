package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/clarkezone/hookserve/hookserve"
	batchv1 "k8s.io/api/batch/v1"
)

const (
	reponame          = "JEKPREV_REPO"
	repopatname       = "JEKPREV_REPO_PAT"
	webhooksecretname = "JEKPREV_WH_SECRET"
	localdirname      = "JEKPREV_LOCALDIR"
	monitorcmdname    = "JEKPREV_monitorCmd"
	initialbranchname = "JEKPREV_initialBranchName"
)

var (
	lrm              *localRepoManager
	jm               *jobmanager
	enableBranchMode bool
)

type cleanupfunc func()

var serve bool
var initialbuild bool
var webhooklisten bool
var initialclone bool
var incluster bool
var sharemgn *httpShareManager

func main() {
	enableBranchMode = false

	// Read and verify flags
	flag.BoolVar(&serve, "serve", false, "start fileserver")
	flag.BoolVar(&initialbuild, "initialbuild", false, "Run an initial build after clone")
	flag.BoolVar(&webhooklisten, "webhooklisten", false, "listen for webhook messages")
	flag.BoolVar(&initialclone, "initialclone", false, "clone repo")
	flag.BoolVar(&incluster, "incluster", false, "Conntect to in-cluster k8s context")
	flag.Parse()

	//repo, repopat, localRootDir, secret, _ := readEnv()
	repo, _, localRootDir, _, _, initalBranchName := readEnv()

	log.Printf("Called with\nrepo:%v\nlocalRootDir:%v\ninitialclone:%v\nwebhooklisten:%v\ninitialbuild:%v\nincluster:%v\nserve:%v\n",
		repo, localRootDir,
		initialclone, webhooklisten, initialbuild, incluster, serve)

	//TODO pass all globals into performactions
	err := PerformActions(repo, localRootDir, initalBranchName, incluster)
	if err != nil {
		log.Printf("Error: %v", err)
		//os.Exit(1)
	}

	// if performactions started the job manager, wait for user to ctrl c out of process
	if jm != nil {
		log.Printf("JobManager exists, initiate wait for interrupt\n")
		//TODO verify this is called when running in cluster
		ch := make(chan struct{})
		handleSig(func() { close(ch) })
		log.Printf("Waiting for user to press control c or sig terminate\n")
		<-ch
		log.Printf("Terminate signal detected, closing job manager\n")
		jm.close()
		log.Printf("Job manager returned from close\n")
		//TODO ? do we need to wait for JM to exit?
		//<-cleanupDone
	}
}

func PerformActions(repo string, localRootDir string, initialBranch string, preformInCluster bool) error {
	if serve || initialbuild || webhooklisten || initialclone {
		result := verifyFlags(repo, localRootDir, initialbuild, initialclone)
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
		err := lrm.initialClone(repo, "")
		if err != nil {
			return err
		}

		if initialBranch != "" {
			return lrm.switchBranch(initialBranch)
		}

	}

	if initialbuild {
		//TODO remove global variable
		jobman, err := newjobmanager(preformInCluster)
		if err != nil {
			return err
		}
		jm = jobman
		notifier := (func(job *batchv1.Job, typee ResourseStateType) {
			log.Printf("Got job in outside world %v", typee)

			if typee == Update && job.Status.Active == 0 && job.Status.Failed > 0 {
				log.Printf("Failed job detected")
			}
		})
		command := []string{"sh", "-c", "--"}
		params := []string{"cd source;bundle install;bundle exec jekyll build -d /site JEKYLL_ENV=production"}
		_, _ = jm.CreateJob("jekyll-render-container", "registry.dev.clarkezone.dev/jekyllbuilder:arm", command, params, notifier)

	}
	return nil
}

func verifyFlags(repo string, localRootDir string, build bool, clone bool) error {
	return nil
	if clone && repo == "" {
		return fmt.Errorf("repo must be provided in %v", reponame)
	}

	if clone {
		if localRootDir == "" {
			return fmt.Errorf("localdir be provided in %v", localRootDir)
		} else {
			fileinfo, res := os.Stat(localRootDir)
			if res != nil {
				return fmt.Errorf("localdir must exist %v", localRootDir)
			}
			if !fileinfo.IsDir() {
				return fmt.Errorf("localdir must be a directory %v", localRootDir)
			}
		}
	}
	if build && !clone {
		return fmt.Errorf("cannont request initial build without an initial clone %v", reponame)
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
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Printf("\nhandleSig Received an interrupt, stopping services...\n")
		if cleanupwork != nil {
			cleanupwork()
		}

		close(cleanupDone)
	}()
	return cleanupDone
}

func readEnv() (string, string, string, string, string, string) {
	repo := os.Getenv(reponame)
	repopat := os.Getenv(repopatname)
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
			lrm.handleWebhook(event.Branch, initialbuild, initialbuild)
		}
	}()
}
