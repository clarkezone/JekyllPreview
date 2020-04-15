package main

import "net/http"

type httpShareManager struct {
	shares map[string]string
}

func createShareManager() *httpShareManager {
	httpMan := &httpShareManager{}
	httpMan.shares = make(map[string]string)
	httpMan.shares["master"] = "master"
	return httpMan
}

func (man *httpShareManager) start() {
	http.Handle("master", http.FileServer(http.Dir("/srv/jekyll/master")))
	http.ListenAndServe(":8085", nil)
}

func (man *httpShareManager) shareBranch(branchName string) {
	_, ok := man.shares[branchName]

	if !ok {
		http.Handle(branchName, http.FileServer(http.Dir(lrm.getCurrentBranchRenderDir())))
		man.shares[branchName] = lrm.getCurrentBranchRenderDir()
	}
}
