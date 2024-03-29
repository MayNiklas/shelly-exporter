package shelly_exporter

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port   string
	listen string

	port_flag   = flag.String("port", "8080", "The port to listen on for HTTP requests.")
	listen_flag = flag.String("listen", "localhost", "The address to listen on for HTTP requests.")
)

func Run() {
	flag.Parse()

	// getting the port from the environment simplifies running the exporter in a container
	// for compatibility with the current NixOS module, we also check the command line flags

	// check if port is set via environment variable
	if port_env := os.Getenv("port"); port_env != "" {
		log.Println("Using port from environment variable:", port_env)
		port = port_env
	} else {
		port = *port_flag
	}

	// check if listen address is set via environment variable
	if listen_env := os.Getenv("listen"); listen_env != "" {
		log.Println("Using listen address from environment variable:", listen_env)
		listen = listen_env
	} else {
		listen = *listen_flag
	}

	log.Println("Starting Shelly exporter on http://" + listen + ":" + port + " ...")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		probeHandler(w, req)
	})

	log.Fatal(http.ListenAndServe(listen+":"+port, nil))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {

	var (
		shelly_power_current = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "shelly_power_current",
				Help: "Current power consumption of shelly.",
			},
			[]string{"name", "hostname", "ip"},
		)
		shelly_power_total = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "shelly_power_total",
				Help: "Total power consumption of shelly.",
			},
			[]string{"name", "hostname", "ip"},
		)
		shelly_uptime = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "shelly_uptime",
				Help: "Uptime of shelly.",
			},
			[]string{"name", "hostname", "ip"},
		)
		shelly_temperature = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "shelly_temperature",
				Help: "Temperature of shelly.",
			},
			[]string{"name", "hostname", "ip"},
		)
		shelly_update_available = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "shelly_update_available",
				Help: "OTA update is available.",
			},
			[]string{"name", "hostname", "ip"},
		)
	)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	r = r.WithContext(ctx)

	// get ?target=<ip> parameter from request
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	// create registry containing metrics
	registry := prometheus.NewPedanticRegistry()

	// add metrics to registry
	registry.MustRegister(shelly_power_current)
	registry.MustRegister(shelly_power_total)
	registry.MustRegister(shelly_temperature)
	registry.MustRegister(shelly_update_available)
	registry.MustRegister(shelly_uptime)

	// get shelly data from target
	var data ShellyData
	if err := data.Fetch(target); err != nil {
		// TODO better error handling
		log.Println(err)
		return
	}

	// define labels used for all metrics
	var shelly_labels prometheus.Labels = prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}

	// set metrics
	shelly_power_current.With(prometheus.Labels(shelly_labels)).Set(data.Status.Meters[0].Power)
	shelly_power_total.With(prometheus.Labels(shelly_labels)).Set(float64(data.Status.Meters[0].Total))
	shelly_temperature.With(prometheus.Labels(shelly_labels)).Set(data.Status.Temperature)
	shelly_uptime.With(prometheus.Labels(shelly_labels)).Set(float64(data.Status.Uptime))

	// check if update is available
	if data.Status.Update.HasUpdate {
		shelly_update_available.With(prometheus.Labels(shelly_labels)).Set(1)
	} else {
		shelly_update_available.With(prometheus.Labels(shelly_labels)).Set(0)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

}
