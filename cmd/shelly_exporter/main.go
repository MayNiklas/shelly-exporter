package shelly_exporter

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type shelly_data struct {
	WifiSta struct {
		Connected bool   `json:"connected"`
		Ssid      string `json:"ssid"`
		IP        string `json:"ip"`
		Rssi      int    `json:"rssi"`
	} `json:"wifi_sta"`
	Cloud struct {
		Enabled   bool `json:"enabled"`
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Mqtt struct {
		Connected bool `json:"connected"`
	} `json:"mqtt"`
	Time          string `json:"time"`
	Unixtime      int    `json:"unixtime"`
	Serial        int    `json:"serial"`
	HasUpdate     bool   `json:"has_update"`
	Mac           string `json:"mac"`
	CfgChangedCnt int    `json:"cfg_changed_cnt"`
	ActionsStats  struct {
		Skipped int `json:"skipped"`
	} `json:"actions_stats"`
	Relays []struct {
		Ison           bool   `json:"ison"`
		HasTimer       bool   `json:"has_timer"`
		TimerStarted   int    `json:"timer_started"`
		TimerDuration  int    `json:"timer_duration"`
		TimerRemaining int    `json:"timer_remaining"`
		Overpower      bool   `json:"overpower"`
		Source         string `json:"source"`
	} `json:"relays"`
	Meters []struct {
		Power     float64   `json:"power"`
		Overpower float64   `json:"overpower"`
		IsValid   bool      `json:"is_valid"`
		Timestamp int       `json:"timestamp"`
		Counters  []float64 `json:"counters"`
		Total     int       `json:"total"`
	} `json:"meters"`
	Temperature     float64 `json:"temperature"`
	Overtemperature bool    `json:"overtemperature"`
	Tmp             struct {
		TC      float64 `json:"tC"`
		TF      float64 `json:"tF"`
		IsValid bool    `json:"is_valid"`
	} `json:"tmp"`
	Update struct {
		Status     string `json:"status"`
		HasUpdate  bool   `json:"has_update"`
		NewVersion string `json:"new_version"`
		OldVersion string `json:"old_version"`
	} `json:"update"`
	RAMTotal int `json:"ram_total"`
	RAMFree  int `json:"ram_free"`
	FsSize   int `json:"fs_size"`
	FsFree   int `json:"fs_free"`
	Uptime   int `json:"uptime"`
}

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func Run() {
	fmt.Println("Starting Shelly exporter!")

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		probeHandler(w, req)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}

// getShellyData returns the data from the shelly device
func getShellyData(shelly_ip string) shelly_data {
	resp, err := http.Get("http://" + shelly_ip + "/status")

	if err != nil {
		fmt.Println("No response from request")
	}

	defer resp.Body.Close()

	// response body is []byte
	body, err := ioutil.ReadAll(resp.Body)
	var result shelly_data

	// Parse []byte to the go struct pointer
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func probeHandler(w http.ResponseWriter, r *http.Request) {

	var (
		shelly_power_current = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "shelly_power_current",
			Help: "Current power consumption of shelly.",
		})
		shelly_power_total = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "shelly_power_total",
			Help: "Total power consumption of shelly.",
		})
		shelly_uptime = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "shelly_uptime",
			Help: "Uptime of shelly.",
		})
		shelly_temperature = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "shelly_temperature",
			Help: "Temperature of shelly.",
		})
		shelly_update_available = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "shelly_update_available",
			Help: "OTA update is available.",
		})
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
	fmt.Println("Probing: ", target)
	var data shelly_data = getShellyData(target)

	// set metrics

	shelly_power_current.Set(data.Meters[0].Power)

	// // shelly_power_total.Set(result.Meters[0].Total) did not work
	// // I'm not 100% sure of the implications of this, but it seems to work.
	shelly_power_total.Set(float64(data.Meters[0].Total))

	shelly_temperature.Set(data.Temperature)

	// check if update is available
	if data.Update.HasUpdate {
		shelly_update_available.Set(1)
	} else {
		shelly_update_available.Set(0)
	}

	shelly_uptime.Set(float64(data.Uptime))

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)

}
