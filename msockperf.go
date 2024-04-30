package main

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func RunSockperf(host string, port string) (string, error) {
	cmd := exec.Command("sockperf", "ping-pong", "-i", host, "-p", port, "-m", "1024", "--tcp")

	fmt.Println(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(output))
	}
	fmt.Println(string(output))
	return string(output), nil
}

type MSockperfObservations struct {
	p25000           float64
	p50000           float64
	p75000           float64
	p90000           float64
	p99000           float64
	p99900           float64
	p99990           float64
	p99999           float64
	runTime          float64
	sentMessages     float64
	receivedMessages float64
	avgLatency       float64
	droppedMessages  float64
}

func MSockGather(host string, port string) MSockperfObservations {
	observations := MSockperfObservations{}

	fmt.Println("Collecting Metrics from Sockperf")
	output, err := RunSockperf(host, port)
	if err != nil {
		return observations
	}

	// Regular expressions to extract different pieces of information
	percentileRegexes := map[string]*regexp.Regexp{
		"99.999": regexp.MustCompile(`percentile 99\.999 =\s+([\d.]+)`),
		"99.990": regexp.MustCompile(`percentile 99\.990 =\s+([\d.]+)`),
		"99.900": regexp.MustCompile(`percentile 99\.900 =\s+([\d.]+)`),
		"99.000": regexp.MustCompile(`percentile 99\.000 =\s+([\d.]+)`),
		"90.000": regexp.MustCompile(`percentile 90\.000 =\s+([\d.]+)`),
		"75.000": regexp.MustCompile(`percentile 75\.000 =\s+([\d.]+)`),
		"50.000": regexp.MustCompile(`percentile 50\.000 =\s+([\d.]+)`),
		"25.000": regexp.MustCompile(`percentile 25\.000 =\s+([\d.]+)`),
	}

	runTimeRegex := regexp.MustCompile(`RunTime=([\d.]+) sec;`)                 // gauge seconds
	sentMessagesRegex := regexp.MustCompile(`SentMessages=([\d]+);`)            // gauge number of messages sent
	receivedMessagesRegex := regexp.MustCompile(`ReceivedMessages=([\d]+)`)     // gauge number of messages received
	latencyRegex := regexp.MustCompile(`Summary: Latency is ([\d.]+)`)          // gauge avg latency in usec
	droppedMessagesRegex := regexp.MustCompile(`# dropped messages = ([\d.]+)`) // counter number of dropped messages

	runTimeMatch := runTimeRegex.FindStringSubmatch(output)

	if len(runTimeMatch) > 1 {
		observations.runTime, _ = strconv.ParseFloat(runTimeMatch[1], 64)
	}

	sentMessagesMatch := sentMessagesRegex.FindStringSubmatch(output)
	if len(sentMessagesMatch) > 1 {
		observations.sentMessages, _ = strconv.ParseFloat(sentMessagesMatch[1], 64)
	}

	receivedMessagesMatch := receivedMessagesRegex.FindStringSubmatch(output)
	if len(receivedMessagesMatch) > 1 {
		observations.receivedMessages, _ = strconv.ParseFloat(receivedMessagesMatch[1], 64)
	}

	latencyMatch := latencyRegex.FindStringSubmatch(output)
	if len(latencyMatch) > 1 {
		observations.avgLatency, _ = strconv.ParseFloat(latencyMatch[1], 64)
	}

	droppedMatch := droppedMessagesRegex.FindStringSubmatch(output)
	if len(droppedMatch) > 1 {
		observations.avgLatency, _ = strconv.ParseFloat(droppedMatch[1], 64)
	}

	for percentile, re := range percentileRegexes {
		match := re.FindStringSubmatch(output)
		match[1] = strings.ReplaceAll(match[1], " ", "")

		if len(match) > 1 {
			switch percentile {

			case "99.999":
				observations.p99999, _ = strconv.ParseFloat(match[1], 64)
			case "99.990":
				observations.p99990, _ = strconv.ParseFloat(match[1], 64)
			case "99.900":
				observations.p99900, _ = strconv.ParseFloat(match[1], 64)
			case "99.000":
				observations.p99000, _ = strconv.ParseFloat(match[1], 64)
			case "90.000":
				observations.p90000, _ = strconv.ParseFloat(match[1], 64)
			case "75.000":
				observations.p75000, _ = strconv.ParseFloat(match[1], 64)
			case "50.000":
				observations.p50000, _ = strconv.ParseFloat(match[1], 64)
			case "25.000":
				observations.p25000, _ = strconv.ParseFloat(match[1], 64)
			}
		}
	}
	return observations
}

func (o *MSockperfObservations) AdjustPercentiles() {
	// Define a slice to hold the percentile values
	percentiles := []float64{o.p25000, o.p50000, o.p75000, o.p90000, o.p99000, o.p99900, o.p99990, o.p99999}

	// Iterate over the percentiles to adjust them
	for i := 0; i < len(percentiles)-1; i++ {
		// Check if current and next percentiles have the same value
		if percentiles[i] >= percentiles[i+1] {
			// Add a small increment to the next percentile value
			percentiles[i+1] = percentiles[i] + 0.0001
		}
	}

	// Update the percentiles in the struct
	o.p25000 = percentiles[0]
	o.p50000 = percentiles[1]
	o.p75000 = percentiles[2]
	o.p90000 = percentiles[3]
	o.p99000 = percentiles[4]
	o.p99900 = percentiles[5]
	o.p99990 = percentiles[6]
	o.p99999 = percentiles[7]
}
