package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(5 * time.Second)
		}
	}()
}

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	shelly_power_current = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_power_current",
		Help: "Current power consumption of shelly.",
	})

	shelly_power_total = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "shelly_power_total",
		Help: "Total power consumption of shelly.",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(shelly_power_current)
	prometheus.MustRegister(shelly_power_total)
}

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Hello world from new Go Collector!")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
