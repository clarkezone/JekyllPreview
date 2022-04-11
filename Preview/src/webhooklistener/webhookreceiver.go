package webhooklistener

import (
	"fmt"

	"github.com/clarkezone/hookserve/hookserve"
	lrm "temp.com/JekyllBlogPreview/localrepomanager"
)

type WebhookListener struct {
	lrm          *lrm.LocalRepoManager
	initialBuild bool
}

func CreateWebhookListener() *WebhookListener {
	wl := WebhookListener{}
	return &wl
}

//nolint
//lint:ignore U1000 called commented out
func (w *WebhookListener) startWebhookListener(secret string) {
	go func() {
		fmt.Printf("Monitoring started\n")
		server := hookserve.NewServer()
		server.Port = 8090
		server.Secret = secret
		server.GoListenAndServe()

		// TODO: make own server with metrics using this for webhooks
		//		hookserve.NewServer().ServeHTTP()

		for event := range server.Events {
			fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
			w.lrm.HandleWebhook(event.Branch, w.initialBuild, w.initialBuild)
		}
	}()
}
