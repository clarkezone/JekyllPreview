package localrepomanager

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"

	batchv1 "k8s.io/api/batch/v1"
	"temp.com/JekyllBlogPreview/jobmanager"
)

const (
	reponame          = "JEKPREV_REPO"
	repopatname       = "JEKPREV_REPO_PAT"
	webhooksecretname = "JEKPREV_WH_SECRET"
	localdirname      = "JEKPREV_LOCALDIR"
	monitorcmdname    = "JEKPREV_monitorCmd"
	initialbranchname = "JEKPREV_initialBranchName"
)

type newBranchHandler interface {
	NewBranch(branch string, dir string)
}

type LocalRepoManager struct {
	currentBranch    string
	repoSourceDir    string
	localRootDir     string
	repo             *gitlayer
	newBranchObs     newBranchHandler
	enableBranchMode bool
	jm               *jobmanager.Jobmanager
}

func CreateLocalRepoManager(rootDir string, newBranch newBranchHandler, enableBranchMode bool, jm *jobmanager.Jobmanager) *LocalRepoManager {
	var lrm = &LocalRepoManager{currentBranch: "master", localRootDir: rootDir}
	lrm.newBranchObs = newBranch
	lrm.enableBranchMode = enableBranchMode
	lrm.jm = jm
	os.RemoveAll(rootDir) // ignore error since it may not exist
	lrm.repoSourceDir = lrm.ensureDir("source")
	return lrm
}

func (lrm *LocalRepoManager) ensureDir(subDir string) string {
	var currentPath = path.Join(lrm.localRootDir, subDir)
	var _, err = os.Stat(currentPath)
	if err != nil {
		err = os.MkdirAll(currentPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Couldn't create sourceDir: %v", err.Error())
		}
	}

	return currentPath
}

func (lrm *LocalRepoManager) getSourceDir() string {
	return lrm.repoSourceDir
}

func (lrm *LocalRepoManager) getRenderDir() string {
	if lrm.enableBranchMode {
		branchName := lrm.legalizeBranchName(lrm.currentBranch)
		return lrm.ensureDir(branchName)
	}
	return lrm.ensureDir("output")
}

func (lrm *LocalRepoManager) legalizeBranchName(name string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(name, "")
}

func (lrm *LocalRepoManager) InitialClone(repo string, repopat string) error {
	//TODO: this function should ensure branch name is correct
	fmt.Printf("Initial clone for\n repo: %v\n local dir:%v", repo, lrm.repoSourceDir)
	if repopat != "" {
		fmt.Printf(" with Personal Access Token.\n")
	} else {
		fmt.Printf(" with no authentication.\n")
	}

	re, err := clone(repo, lrm.repoSourceDir, repopat)
	if err != nil {
		fmt.Printf("Error in initial clone: %v\n", err.Error())
		os.Exit(1)
	}
	lrm.repo = re
	fmt.Println("Clone Done.")
	return err
}

func (lrm *LocalRepoManager) SwitchBranch(branch string) error {
	if branch != lrm.currentBranch {
		fmt.Printf("Fetching\n")

		err := lrm.repo.checkout(branch)
		if err != nil {
			return fmt.Errorf("checkout failed: %v", err.Error())
		}

		lrm.currentBranch = branch
	}

	err := lrm.repo.pull(branch)
	if err != nil {
		return fmt.Errorf("pull failed: %v", err.Error())
	}
	return nil
}

//nolint
//lint:ignore U1000 called commented out
func (lrm *LocalRepoManager) HandleWebhook(branch string, runjek bool, sendNotify bool) {
	err := lrm.SwitchBranch(branch)
	if err != nil {
		panic(err)
	}

	renderDir := lrm.getRenderDir()
	// todo handle branch change
	lrm.StartJob()

	if lrm.enableBranchMode && sendNotify && lrm.newBranchObs != nil {
		lrm.newBranchObs.NewBranch(lrm.legalizeBranchName(branch), renderDir)
	}
}

func (lrm *LocalRepoManager) StartJob() {
	if lrm.jm == nil {
		log.Println("Skipping StartJob due to lack of jobmanager instance")
		return
	}
	namespace := "jekyllpreviewv2"
	notifier := (func(job *batchv1.Job, typee jobmanager.ResourseStateType) {
		log.Printf("Got job in outside world %v", typee)

		if typee == jobmanager.Update && job.Status.Active == 0 && job.Status.Failed > 0 {
			log.Printf("Failed job detected")
		}
	})
	var imagePath string
	fmt.Printf("%v", runtime.GOARCH)
	if runtime.GOARCH == "amd64" {
		imagePath = "registry.hub.docker.com/clarkezone/jekyllbuilder:0.0.1.8"
	} else {

		imagePath = "registry.dev.clarkezone.dev/jekyllbuilder:arm"
	}
	command := []string{"sh", "-c", "--"}
	params := []string{"cd source;bundle install;bundle exec jekyll build -d /site JEKYLL_ENV=production"}
	_, err := lrm.jm.CreateJob("jekyll-render-container", namespace, imagePath, command, params, notifier)
	if err != nil {
		log.Printf("Failed to create job: %v\n", err.Error())
	}
}

func ReadEnv() (string, string, string, string, string, string) {
	repo := os.Getenv(reponame)
	if repo == "" {
		err := fmt.Sprintf("Environment variable %v was empty", reponame)
		panic(err)
	}
	repopat := os.Getenv(repopatname)
	localdr := os.Getenv(localdirname)
	secret := os.Getenv(webhooksecretname)
	monitorcmdline := os.Getenv(monitorcmdname)
	initalbranchname := os.Getenv(initialbranchname)
	return repo, repopat, localdr, secret, monitorcmdline, initalbranchname
}
