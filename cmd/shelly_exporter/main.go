package shelly_exporter

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port   = flag.String("port", "8080", "The port to listen on for HTTP requests.")
	listen = flag.String("listen", "localhost", "The address to listen on for HTTP requests.")
)

func Run() {
	flag.Parse()

	log.Println("Starting Shelly exporter on http://" + *listen + ":" + *port + " ...")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		probeHandler(w, req)
	})

	log.Fatal(http.ListenAndServe(*listen+":"+*port, nil))
}

func probeHandler(w http.ResponseWriter, r *http.Request) {

	var (
		shelly_power_current = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "shelly_power_current",
				Help: "Current power consumption of shelly.",
			})
		shelly_power_total = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "shelly_power_total",
				Help: "Total power consumption of shelly.",
			})
		shelly_uptime = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "shelly_uptime",
				Help: "Uptime of shelly.",
			})
		shelly_temperature = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "shelly_temperature",
				Help: "Temperature of shelly.",
			})
		shelly_update_available = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "shelly_update_available",
				Help: "OTA update is available.",
			})
		shelly_name = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "shelly_name",
				Help: "Name of shelly.",
			},
			[]string{"name", "hostname"},
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
	registry.MustRegister(shelly_name)
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

	// set metrics
	shelly_name.With(prometheus.Labels{"name": data.Name, "hostname": data.Device.Hostname})
	shelly_power_current.Set(data.Meters[0].Power)
	shelly_power_total.Set(float64(data.Meters[0].Total))
	shelly_temperature.Set(data.Temperature)
	shelly_uptime.Set(float64(data.Uptime))

	// check if update is available
	if data.Update.HasUpdate {
		shelly_update_available.Set(1)
	} else {
		shelly_update_available.Set(0)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

}
