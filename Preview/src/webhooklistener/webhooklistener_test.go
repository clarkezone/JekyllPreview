package webhooklistener

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	hs "github.com/clarkezone/hookserve/hookserve"
)

func GetBody() *strings.Reader {
	body := `{
		"ref": "refs/heads/master",
		"before": "cfbda0818c286970d373bab5e700599e572d3c40",
		"after": "e2cdd9288b113800293027bc5dae2c7d47b36189",
		"compare_url": "http://gitea.homelab.clarkezone.dev:3000/clarkezone/testfoobar2/compare/cfbda0818c286970d373bab5e700599e572d3c40...e2cdd9288b113800293027bc5dae2c7d47b36189",
		"commits": [
		  {
			"id": "e2cdd9288b113800293027bc5dae2c7d47b36189",
			"message": "Update 'test.txt'\n",
			"url": "http://gitea.homelab.clarkezone.dev:3000/clarkezone/testfoobar2/commit/e2cdd9288b113800293027bc5dae2c7d47b36189",
			"author": {
			  "name": "clarkezone",
			  "email": "james@clarkezone.io",
			  "username": "clarkezone"
			},
			"committer": {
			  "name": "clarkezone",
			  "email": "james@clarkezone.io",
			  "username": "clarkezone"
			},
			"verification": null,
			"timestamp": "2022-04-10T09:35:51Z",
			"added": [],
			"removed": [],
			"modified": [
			  "test.txt"
			]
		  }
		],
		"head_commit": {
		  "id": "e2cdd9288b113800293027bc5dae2c7d47b36189",
		  "message": "Update 'test.txt'\n",
		  "url": "http://gitea.homelab.clarkezone.dev:3000/clarkezone/testfoobar2/commit/e2cdd9288b113800293027bc5dae2c7d47b36189",
		  "author": {
			"name": "clarkezone",
			"email": "james@clarkezone.io",
			"username": "clarkezone"
		  },
		  "committer": {
			"name": "clarkezone",
			"email": "james@clarkezone.io",
			"username": "clarkezone"
		  },
		  "verification": null,
		  "timestamp": "2022-04-10T09:35:51Z",
		  "added": [],
		  "removed": [],
		  "modified": [
			"test.txt"
		  ]
		},
		"repository": {
		  "id": 1,
		  "owner": {"id":1,"login":"clarkezone","full_name":"","email":"james@clarkezone.io","avatar_url":"http://gitea.homelab.clarkezone.dev:3000/user/avatar/clarkezone/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2021-11-21T18:43:19Z","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":0,"username":"clarkezone"},
		  "name": "testfoobar2",
		  "full_name": "clarkezone/testfoobar2",
		  "description": "",
		  "empty": false,
		  "private": false,
		  "fork": false,
		  "template": false,
		  "parent": null,
		  "mirror": false,
		  "size": 21,
		  "html_url": "http://gitea.homelab.clarkezone.dev:3000/clarkezone/testfoobar2",
		  "ssh_url": "ssh://git@gitea.homelab.clarkezone.dev:2222/clarkezone/testfoobar2.git",
		  "clone_url": "http://gitea.homelab.clarkezone.dev:3000/clarkezone/testfoobar2.git",
		  "original_url": "",
		  "website": "",
		  "stars_count": 0,
		  "forks_count": 0,
		  "watchers_count": 1,
		  "open_issues_count": 0,
		  "open_pr_counter": 0,
		  "release_counter": 0,
		  "default_branch": "master",
		  "archived": false,
		  "created_at": "2021-11-21T18:50:53Z",
		  "updated_at": "2022-04-10T09:32:02Z",
		  "permissions": {
			"admin": true,
			"push": true,
			"pull": true
		  },
		  "has_issues": true,
		  "internal_tracker": {
			"enable_time_tracker": true,
			"allow_only_contributors_to_track_time": true,
			"enable_issue_dependencies": true
		  },
		  "has_wiki": true,
		  "has_pull_requests": true,
		  "has_projects": true,
		  "ignore_whitespace_conflicts": false,
		  "allow_merge_commits": true,
		  "allow_rebase": true,
		  "allow_rebase_explicit": true,
		  "allow_squash_merge": true,
		  "default_merge_style": "merge",
		  "avatar_url": "",
		  "internal": false,
		  "mirror_interval": ""
		},
		"pusher": {"id":1,"login":"clarkezone","full_name":"","email":"james@clarkezone.io","avatar_url":"http://gitea.homelab.clarkezone.dev:3000/user/avatar/clarkezone/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2021-11-21T18:43:19Z","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":0,"username":"clarkezone"},
		"sender": {"id":1,"login":"clarkezone","full_name":"","email":"james@clarkezone.io","avatar_url":"http://gitea.homelab.clarkezone.dev:3000/user/avatar/clarkezone/-1","language":"","is_admin":false,"last_login":"0001-01-01T00:00:00Z","created":"2021-11-21T18:43:19Z","restricted":false,"active":false,"prohibit_login":false,"location":"","website":"","description":"","visibility":"public","followers_count":0,"following_count":0,"starred_repos_count":0,"username":"clarkezone"}
	  }
	`
	reader := strings.NewReader(body)
	return reader
}

func TestGiteaParse(t *testing.T) {
	reader := GetBody()
	req := httptest.NewRequest(http.MethodPost, "/postreceive", reader)
	req.Header.Set("X-GitHub-Event", "push")
	w := httptest.NewRecorder()

	serv := hs.NewServer()

	serv.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Request failed due to bad payload")
	}

	event := <-serv.Events
	fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)

	if event.Type != "push" || event.Branch != "master" || event.Repo != "testfoobar2" {
		t.Errorf("didn't match.")
	}
}

func Test_webhooklistening(t *testing.T) {
	//wait := make(chan bool)
	wh := WebhookListener{}
	wh.StartListen("ss")
	reader := GetBody()
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://0.0.0.0:8090", reader)
	if err != nil {
		t.Errorf("bad request")
	}
	req.Header.Set("X-GitHub-Event", "push")
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//<-wait
	err = wh.Shutdown()
	if err != nil {
		t.Errorf("shutdown failed")
	}
}
