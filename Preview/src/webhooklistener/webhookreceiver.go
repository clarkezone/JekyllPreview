package webhooklistener

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

func CreateWebhookListener(lrm *lrm.LocalRepoManager) *WebhookListener {
	wl := WebhookListener{}
	wl.lrm = lrm
	return &wl
}

func (wl *WebhookListener) StartListen(secret string) {
	fmt.Println("starting...")

	//prometheus.MustRegister(requestDuration)

	wl.hookserver = hookserve.NewServer()
	wl.ctx, wl.cancel = context.WithCancel(context.Background())
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
		err := wl.httpserver.ListenAndServe()
		if err != nil {
			panic(err)
		}
		defer func() {
			log.Println("Webserver exited")
		}()
	}()
	go func() {
		defer func() {
			log.Println("processing loop exited")
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
