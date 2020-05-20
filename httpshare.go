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
		//http.Handle(httpBranchName, http.StripPrefix(httpBranchName, http.FileServer(http.Dir(dir))))
		http.Handle(httpBranchName, auth(handleFileServer(dir, httpBranchName)))
		man.shares[httpBranchName] = dir
	}
}

func (man *httpShareManager) shareRootDir(dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
}

func (man *httpShareManager) NewBranch(branchName string, dir string) {
	man.shareBranch(branchName, dir)
}

func handleFileServer(dir, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		realHandler(w, req)
	}
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if !check(user, pass) {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}

func check(u, p string) bool {
	if u == "" && p == "" {
		return true
	}
	return false
}
