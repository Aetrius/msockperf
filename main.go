package main

import (
	//"errors"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	latencySummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			//Namespace:  namespace,
			Name:       "msockperf_latency_summary",
			Help:       "Summary of sockperf latency in [Microseconds]",
			Objectives: map[float64]float64{0.25: 0.25, 0.5: 0.05, 0.75: 0.05, 0.9: 0.9, 0.99: 0.99, 0.999: 0.999, 0.9999: 0.9999},
		},
		[]string{"namespace"},
	)

	runTimeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "msockperf_runtime_gauge",
			Help: "Total Runtime of the sockperf benchmark test - [seconds]",
		},
		[]string{"namespace"}, // Define the namespace label here
	)

	sentMessagesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_sent_messages_gauge",
			Help: "Total messages sent through the sockperf benchmark test - [quantity of messages]",
		},
		[]string{"namespace"}, // Define the namespace label here
	)

	receivedMessagesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_received_messages_gauge",
			Help: "Total messages received through the sockperf benchmark test - [quantity of messages]",
		},
		[]string{"namespace"}, // Define the namespace label here
	)

	latencyAvgGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_avg_latency_gauge",
			Help: "Average Latency of the msockperf benchmark test - [usec - Microseconds]",
		},
		[]string{"namespace"}, // Define the namespace label here
	)

	droppedMessagesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_dropped_messages_gauge",
			Help: "Count of dropped messages from the msockperf benchmark test - [quantity of messages]",
		},
		[]string{"namespace"}, // Define the namespace label here
	)
)

// Function to get the value of an environment variable
func getEnvVars(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		// If the environment variable is not set, return the default value
		return defaultValue
	}
	return value
}

func main() {
	rand.Seed(time.Now().Unix())
	// os package

	// Default values
	host := "127.0.0.1"
	port := "11111"
	namespace := "default"

	// Pull values for the run
	host = getEnvVars("MSOCK_REMOTE_HOST", host)
	port = getEnvVars("MSOCK_REMOTE_PORT", port)
	namespace = getEnvVars("MSOCK_NAMESPACE", namespace)

	// Resolve host if it's not an IP address
	host = resolveHost(host)

	fmt.Println("Host: ", host)
	fmt.Println("Port: ", port)
	fmt.Println("Namespace: ", namespace)

	prometheus.Register(latencySummary)
	prometheus.Register(runTimeGauge)
	prometheus.Register(sentMessagesGauge)
	prometheus.Register(receivedMessagesGauge)
	prometheus.Register(latencyAvgGauge)
	prometheus.Register(droppedMessagesGauge)

	http.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), latencySummary, host, port, namespace))

	// go func() {
	// 	for {
	// 		//observations := MSockGather(host, port)
	// 		//observations.AdjustPercentiles()
	// 		// Sleep to avoid busy looping
	// 		//time.Sleep(time.Second)
	// 	}
	// }()

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func newHandlerWithHistogram(handler http.Handler, summary *prometheus.SummaryVec, host string, port string, namespace string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		status := http.StatusOK
		defer func() {
			observations := MSockGather(host, port)
			observations.AdjustPercentiles()

			// Observe values in the summary
			summary.WithLabelValues(namespace).Observe(observations.p25000)
			summary.WithLabelValues(namespace).Observe(observations.p50000)
			summary.WithLabelValues(namespace).Observe(observations.p75000)
			summary.WithLabelValues(namespace).Observe(observations.p90000)
			summary.WithLabelValues(namespace).Observe(observations.p99000)
			summary.WithLabelValues(namespace).Observe(observations.p99900)
			summary.WithLabelValues(namespace).Observe(observations.p99990)

			//runTimeGauge.Add(observations.runTime)
			//sentMessagesGauge.Add(observations.sentMessages)
			//receivedMessagesGauge.Add(observations.receivedMessages)
			//latencyAvgGauge.Add(observations.avgLatency)
			//droppedMessagesGauge.Add(observations.droppedMessages)
			runTimeGauge.WithLabelValues(namespace).Set(observations.runTime)
			sentMessagesGauge.WithLabelValues(namespace).Set(observations.sentMessages)
			receivedMessagesGauge.WithLabelValues(namespace).Set(observations.receivedMessages)
			latencyAvgGauge.WithLabelValues(namespace).Set(observations.avgLatency)
			droppedMessagesGauge.WithLabelValues(namespace).Set(observations.droppedMessages)

		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}
