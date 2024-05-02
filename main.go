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
		[]string{"namespace", "pod_ip", "node_name", "destination"},
	)

	runTimeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "msockperf_runtime_gauge",
			Help: "Total Runtime of the sockperf benchmark test - [seconds]",
		},
		[]string{"namespace", "pod_ip", "node_name", "destination"}, // Define the namespace label here
	)

	sentMessagesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_sent_messages_gauge",
			Help: "Total messages sent through the sockperf benchmark test - [quantity of messages]",
		},
		[]string{"namespace", "pod_ip", "node_name", "destination"}, // Define the namespace label here
	)

	receivedMessagesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_received_messages_gauge",
			Help: "Total messages received through the sockperf benchmark test - [quantity of messages]",
		},
		[]string{"namespace", "pod_ip", "node_name", "destination"}, // Define the namespace label here
	)

	latencyAvgGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_avg_latency_gauge",
			Help: "Average Latency of the msockperf benchmark test - [usec - Microseconds]",
		},
		[]string{"namespace", "pod_ip", "node_name", "destination"}, // Define the namespace label here
	)

	droppedMessagesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			//Namespace: "default",
			Name: "msockperf_dropped_messages_gauge",
			Help: "Count of dropped messages from the msockperf benchmark test - [quantity of messages]",
		},
		[]string{"namespace", "pod_ip", "node_name", "destination"}, // Define the namespace label here
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

	// Default values
	host := "127.0.0.1"
	port := "11111"
	namespace := "default"
	podIp := "undefined"
	nodeName := "undefined"
	node := "undefined"

	// Pull values for the run
	host = getEnvVars("MSOCK_REMOTE_HOST", host)
	port = getEnvVars("MSOCK_REMOTE_PORT", port)
	namespace = getEnvVars("MSOCK_NAMESPACE", namespace)
	podIp = getEnvVars("POD_IP", podIp)
	nodeName = getEnvVars("NODE_NAME", nodeName)
	node = getEnvVars("NODE", nodeName)

	// Resolve host if it's not an IP address
	host = resolveHost(host)

	if nodeName == "undefined" {
		nodeName = node
	}

	fmt.Println("Host: ", host)
	fmt.Println("Port: ", port)
	fmt.Println("Namespace: ", namespace)
	fmt.Println("Pod IP: ", podIp)
	fmt.Println("Node Name ", nodeName)
	fmt.Println("Node ", node)

	prometheus.Register(latencySummary)
	prometheus.Register(runTimeGauge)
	prometheus.Register(sentMessagesGauge)
	prometheus.Register(receivedMessagesGauge)
	prometheus.Register(latencyAvgGauge)
	prometheus.Register(droppedMessagesGauge)

	http.Handle("/metrics", newHandlerWithHistogram(promhttp.Handler(), latencySummary, host, port, namespace, podIp, nodeName))

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

func newHandlerWithHistogram(handler http.Handler, summary *prometheus.SummaryVec, host string, port string, namespace string, podIp string, nodeName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		status := http.StatusOK
		defer func() {
			observations := MSockGather(host, port)
			observations.AdjustPercentiles()

			// Observe values in the summary
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p25000)
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p50000)
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p75000)
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p90000)
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p99000)
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p99900)
			summary.WithLabelValues(namespace, podIp, nodeName, host).Observe(observations.p99990)

			runTimeGauge.WithLabelValues(namespace, podIp, nodeName, host).Set(observations.runTime)
			sentMessagesGauge.WithLabelValues(namespace, podIp, nodeName, host).Set(observations.sentMessages)
			receivedMessagesGauge.WithLabelValues(namespace, podIp, nodeName, host).Set(observations.receivedMessages)
			latencyAvgGauge.WithLabelValues(namespace, podIp, nodeName, host).Set(observations.avgLatency)
			droppedMessagesGauge.WithLabelValues(namespace, podIp, nodeName, host).Set(observations.droppedMessages)

		}()

		if req.Method == http.MethodGet {
			handler.ServeHTTP(w, req)
			return
		}
		status = http.StatusBadRequest

		w.WriteHeader(status)
	})
}
