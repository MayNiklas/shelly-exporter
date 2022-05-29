package shelly_exporter

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			requestShelly()

			time.Sleep(5 * time.Second)
		}
	}()
}

func requestShelly() {
	// Get request
	resp, err := http.Get("http://192.168.15.2/status")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result shelly_data
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	shelly_power_current.Set(result.Meters[0].Power)

	// shelly_power_total.Set(result.Meters[0].Total) did not work
	// "cannot use result.Meters[0].Total (variable of type int) as float64 value in argument to shelly_power_total.
	// SetcompilerIncompatibleAssign"
	// I'm not 100% sure of the implications of this, but it seems to work.
	shelly_power_total.Set(float64(result.Meters[0].Total))

	shelly_uptime.Set(float64(result.Uptime))

	shelly_temperature.Set(result.Temperature)

	if result.Update.HasUpdate {
		shelly_update_available.Set(1)
	} else {
		shelly_update_available.Set(0)
	}
}

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(shelly_power_current)
	prometheus.MustRegister(shelly_power_total)
	prometheus.MustRegister(shelly_uptime)
	prometheus.MustRegister(shelly_temperature)
	prometheus.MustRegister(shelly_update_available)
}

func Run() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Hello world from new Go Collector!")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
