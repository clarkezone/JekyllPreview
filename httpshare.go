package main

import "net/http"

type httpShareManager struct {
	shares map[string]string
}

func createShareManager() *httpShareManager {
	httpMan := &httpShareManager{}
	httpMan.shares = make(map[string]string)
	return httpMan
}

func (man *httpShareManager) start() {
	go func() { http.ListenAndServe(":8085", nil) }()
}

func (man *httpShareManager) shareBranch(branchName string, dir string) {
	httpBranchName := "/" + branchName + "/"
	_, ok := man.shares[httpBranchName]

	if !ok {
		http.Handle(httpBranchName, http.StripPrefix(httpBranchName, http.FileServer(http.Dir(dir))))
		man.shares[httpBranchName] = dir
	}
}

func (man *httpShareManager) NewBranch(branchName string, dir string) {
	man.shareBranch(branchName, dir)
}
