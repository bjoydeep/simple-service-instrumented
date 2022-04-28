package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var counter int
var mutex = &sync.Mutex{}
var (
	opsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "a_hits_total",
		Help: "The total number of hits",
	}, []string{"method", "route"})
)

var rpcDurationVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "avec_duration_seconds",
	Help:    "RPC latency distributions.",
	Buckets: prometheus.DefBuckets, //buckets can be customized as well
}, []string{"method", "route", "status_code"})

var testData = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "a_test_performance_seconds",
	Help: "99th percentile latency in seconds.",
}, []string{"method", "route", "snapshot"})

var oscillationPeriod = flag.Duration("oscillation-period", 10*time.Minute, "The duration of the rate oscillation period.")

func incrementCounter() {
	mutex.Lock()
	counter++
	mutex.Unlock()
}

func fetchTestData() {

	// main point is that this data can be fetched from elsewhere
	// and not in line code - like the counters
	// and we could make blocking calls here if needed
	testData.WithLabelValues("GET", "/health", "1.0.0").Add(0.009)
	testData.WithLabelValues("GET", "/health", "2.0.0").Add(0.009)
	testData.WithLabelValues("GET", "/test", "1.0.0").Add(2.354)
	testData.WithLabelValues("GET", "/test", "2.0.0").Add(2.354)
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, how are you !!\n")

	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//steps to add a random delay
		oscillationFactor := func() float64 {
			return 2 + math.Sin(math.Sin(2*math.Pi*float64(time.Since(start))/float64(*oscillationPeriod)))
		}
		fmt.Fprintf(w, "Calling health endpoint. Service is up\n")
		//just to check if prom is capturing the correct numbers
		incrementCounter()
		fmt.Fprintf(w, strconv.Itoa(counter))
		opsProcessed.WithLabelValues(r.Method, r.RequestURI).Inc()
		//Adding random delay to make the histogram data interesting
		time.Sleep(time.Duration(100*oscillationFactor()) * time.Millisecond)
		rpcDurationVec.WithLabelValues(r.Method, r.RequestURI, "HTTP Status Code goes here").Observe(time.Since(start).Seconds())

	})


	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
