package main

import (
	"fmt"
	"net/http"
	"strings"
)

type httpShareManager struct {
	shares          map[string]string
	subdomainShares map[string]http.Handler
}

func createShareManager() *httpShareManager {
	httpMan := &httpShareManager{}
	httpMan.shares = make(map[string]string)
	httpMan.subdomainShares = make(map[string]http.Handler)
	return httpMan
}

func (man *httpShareManager) start() {
	go func() { http.ListenAndServe(":8085", nil) }()
}

func (man *httpShareManager) startsubdomain() {
	go func() { http.ListenAndServe(":8085", man) }()
}

func (man *httpShareManager) shareBranchPath(branchName string, dir string) {
	httpBranchName := "/" + branchName + "/"
	_, ok := man.shares[httpBranchName]

	if !ok {
		http.Handle(httpBranchName, http.StripPrefix(httpBranchName, http.FileServer(http.Dir(dir))))
		man.shares[httpBranchName] = dir
	}
}

func (man *httpShareManager) shareBranchSubdomain(branchName string, dir string) {
	branchserver := http.FileServer(http.Dir(dir))
	man.subdomainShares[branchName] = branchserver
}

func (man *httpShareManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO verify that preview is present
	domainParts := strings.Split(r.Host, ".")

	if mux := man.subdomainShares[domainParts[0]]; mux != nil {
		// Let the appropriate mux serve the request
		mux.ServeHTTP(w, r)
	} else {
		// Handle 404
		error := fmt.Sprintf("Found found subdomain:%v: ", domainParts[0])
		http.Error(w, error, 404)
	}
}

func (man *httpShareManager) shareRootDir(dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
}

func (man *httpShareManager) NewBranch(branchName string, dir string) {
	//man.shareBranchPath(branchName, dir)
	man.shareBranchSubdomain(branchName, dir)
}
