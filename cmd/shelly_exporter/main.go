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

	// set metrics
	shelly_power_current.With(prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}).Set(data.Status.Meters[0].Power)
	shelly_power_total.With(prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}).Set(float64(data.Status.Meters[0].Total))
	shelly_temperature.With(prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}).Set(data.Status.Temperature)
	shelly_uptime.With(prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}).Set(float64(data.Status.Uptime))

	// check if update is available
	if data.Status.Update.HasUpdate {
		shelly_update_available.With(prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}).Set(1)
	} else {
		shelly_update_available.With(prometheus.Labels{"name": data.Settings.Name, "hostname": data.Settings.Device.Hostname, "ip": data.Status.WifiSta.IP}).Set(0)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

}
