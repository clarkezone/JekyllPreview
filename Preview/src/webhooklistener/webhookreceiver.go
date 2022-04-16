package webhooklistener

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	lrm "temp.com/JekyllBlogPreview/localrepomanager"

	"github.com/clarkezone/hookserve/hookserve"
)

//var requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
//	Name: "go_request_operations_total",
//	Help: "The total number of processed requests",
//})
//
//var requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
//	Name:    "go_request_duration_seconds",
//	Help:    "Histogram for the duration in seconds.",
//	Buckets: []float64{1, 2, 5, 6, 10},
//},
//	[]string{"endpoint"},
//)

type WebhookListener struct {
	lrm          *lrm.LocalRepoManager
	initialBuild bool
	hookserver   *hookserve.Server
	httpserver   *http.Server
	ctx          context.Context
	cancel       context.CancelFunc
}

func CreateWebhookListener() *WebhookListener {
	wl := WebhookListener{}
	return &wl
}

//nolint
//lint:ignore U1000 called commented out
func (w *WebhookListener) StartWebhookListener(secret string) {
	go func() {
		fmt.Printf("Monitoring started\n")

		server := hookserve.NewServer()
		server.Port = 8090
		server.Secret = secret
		server.GoListenAndServe()

		for event := range server.Events {
			fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
			w.lrm.HandleWebhook(event.Branch, w.initialBuild, w.initialBuild)
		}
	}()
}

func (wl *WebhookListener) StartListen(secret string) {
	fmt.Println("starting...")

	//prometheus.MustRegister(requestDuration)

	wl.hookserver = hookserve.NewServer()
	wl.ctx, wl.cancel = context.WithTimeout(context.Background(), 5*time.Second)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//start a timer
		//start := time.Now()

		//Call webhooklistener
		wl.hookserver.ServeHTTP(w, r)

		//measure the duration and log to prometheus
		//httpDuration := time.Since(start)
		//requestDuration.WithLabelValues("GET /").Observe(httpDuration.Seconds())

		//increment a counter for number of requests processed
		//requestsProcessed.Inc()
	})
	wl.httpserver = &http.Server{Addr: ":8090"}
	//http.Handle("/metrics", promhttp.Handler())
	go func() {
		wl.httpserver.ListenAndServe()
		defer func() {
			log.Println("Exited")
		}()
		for {
			select {
			case <-wl.ctx.Done():
				return
			case event := <-wl.hookserver.Events:
				fmt.Println(event.Owner + " " + event.Repo + " " + event.Branch + " " + event.Commit)
				wl.lrm.HandleWebhook(event.Branch, wl.initialBuild, wl.initialBuild)
			}
		}
	}()
}

func (wl *WebhookListener) Shutdown() error {
	defer wl.ctx.Done()
	defer wl.cancel()
	return wl.httpserver.Shutdown(wl.ctx)
}