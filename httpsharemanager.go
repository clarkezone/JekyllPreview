package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

type httpShareManager struct {
	shares           map[string]string
	subdomainShares  map[string]http.Handler
	topLevelMux      http.Handler
	authHandler      http.Handler
	sslhostwhitelist string
}

func createShareManager(toplevelmux http.Handler, sslhostwhitelist string) *httpShareManager {
	httpMan := &httpShareManager{}
	httpMan.shares = make(map[string]string)
	httpMan.subdomainShares = make(map[string]http.Handler)
	httpMan.topLevelMux = toplevelmux
	httpMan.sslhostwhitelist = sslhostwhitelist
	return httpMan
}

func (man *httpShareManager) start() {
	go func() {
		if man.sslhostwhitelist != "" {
			http.ListenAndServe(":8085", nil)
		} else {
			http.ListenAndServe(":8085", nil)
		}

	}()
}

func (man *httpShareManager) startsubdomain() {
	go func() {
		if man.sslhostwhitelist != "" {
			certManager, config := man.getssl(man.sslhostwhitelist)

			server := &http.Server{
				Addr:      ":8443",
				TLSConfig: config,
				Handler:   man,
			}

			go http.ListenAndServe(":8080", certManager.HTTPHandler(nil))

			server.ListenAndServeTLS("", "")
		} else {
			http.ListenAndServe(":8085", man)
		}

	}()
}

func (man *httpShareManager) getssl(whileList string) (*autocert.Manager, *tls.Config) {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: man.GetPolicy,
		Cache:      autocert.DirCache("certs"),
	}

	config := &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}
	return &certManager, config
}

func (man *httpShareManager) GetPolicy(ctx context.Context, host string) error {
	if !strings.HasSuffix(host, man.sslhostwhitelist) {
		return errors.New("bad host")
	}
	return nil
}

func (man *httpShareManager) shareBranchPath(branchName string, dir string) {
	httpBranchName := "/" + branchName + "/"
	_, ok := man.shares[httpBranchName]

	if !ok {
		http.Handle(httpBranchName, http.StripPrefix(httpBranchName, http.FileServer(http.Dir(dir))))
		man.shares[httpBranchName] = dir
	}
}

func handleFileServer(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix("", fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		realHandler(w, req)
	}
}

func (man *httpShareManager) shareBranchSubdomain(branchName string, dir string) {
	//branchserver := http.FileServer(http.Dir(dir))
	branchserver := auth(handleFileServer(dir))
	man.subdomainShares[branchName] = branchserver
}

func (man *httpShareManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainParts := strings.Split(r.Host, ".")
	if len(domainParts) >= 3 && domainParts[0] != "preview" {
		if mux := man.subdomainShares[domainParts[0]]; mux != nil {
			// Let the appropriate mux serve the request

			mux.ServeHTTP(w, r)
		} else {
			// Handle 404
			error := fmt.Sprintf("Not found :%v: ", domainParts[0])
			http.Error(w, error, 404)
		}
	} else if domainParts[0] == "preview" {
		man.topLevelMux.ServeHTTP(w, r)
	}
}

func (man *httpShareManager) shareRootDir(dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
}

func (man *httpShareManager) NewBranch(branchName string, dir string) {
	//man.shareBranchPath(branchName, dir)
	man.shareBranchSubdomain(branchName, dir)
}
