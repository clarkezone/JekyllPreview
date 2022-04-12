package webhooklistener

import (
	"fmt"
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

func (w *WebhookListener) StartListen(secret string) {
	fmt.Println("starting...")

	//prometheus.MustRegister(requestDuration)

	serv := hookserve.NewServer()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//start a timer
		//start := time.Now()

		//Call webhooklistener
		serv.ServeHTTP(w, r)

		//measure the duration and log to prometheus
		//httpDuration := time.Since(start)
		//requestDuration.WithLabelValues("GET /").Observe(httpDuration.Seconds())

		//increment a counter for number of requests processed
		//requestsProcessed.Inc()
	})

	//http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8090", nil)
}
