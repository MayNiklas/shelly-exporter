package shelly_exporter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type ShellyData struct {
	Status   ShellyStatus
	Settings ShellySettings
	// shelly_status
	// TODO json tag "wifi_sta" is present in both structs and will lead to problems!
	// shelly_settings
}

func getJson(url string) ([]byte, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (s *ShellyData) Fetch(address string) error {

	var (
		statusJson   []byte
		settingsJson []byte
		err          error
	)

	if statusJson, err = getJson(address + "/status"); err != nil {
		return err
	}

	if err := json.Unmarshal(statusJson, &s.Status); err != nil {
		return err
	}

	// Only fetch settings, if Name is unset. Fetching settings once is should
	// be sufficient and save bandwith
	if len(s.Settings.Name) == 0 {

		if settingsJson, err = getJson(address + "/settings"); err != nil {
			return err
		}

		if err := json.Unmarshal(settingsJson, &s.Settings); err != nil {
			return err
		}
	}

	return nil
}

type ShellyStatus struct {
	WifiSta         WifiSta      `json:"wifi_sta"`
	Cloud           Cloud        `json:"cloud"`
	Mqtt            Mqtt         `json:"mqtt"`
	Time            string       `json:"time"`
	Unixtime        int          `json:"unixtime"`
	Serial          int          `json:"serial"`
	HasUpdate       bool         `json:"has_update"`
	Mac             string       `json:"mac"`
	CfgChangedCnt   int          `json:"cfg_changed_cnt"`
	ActionsStats    ActionsStats `json:"actions_stats"`
	Relays          []Relays     `json:"relays"`
	Meters          []Meters     `json:"meters"`
	Temperature     float64      `json:"temperature"`
	Overtemperature bool         `json:"overtemperature"`
	Tmp             Tmp          `json:"tmp"`
	Update          Update       `json:"update"`
	RAMTotal        int          `json:"ram_total"`
	RAMFree         int          `json:"ram_free"`
	FsSize          int          `json:"fs_size"`
	FsFree          int          `json:"fs_free"`
	Uptime          int          `json:"uptime"`
}

type WifiSta struct {
	Enabled    bool   `json:"enabled"`
	Ssid       string `json:"ssid"`
	Ipv4Method string `json:"ipv4_method"`
	Gw         string `json:"gw"`
	Mask       string `json:"mask"`
	DNS        string `json:"dns"`
	Connected  bool   `json:"connected"`
	IP         string `json:"ip"`
	Rssi       int    `json:"rssi"`
}

type Cloud struct {
	Enabled   bool `json:"enabled"`
	Connected bool `json:"connected"`
}

type ActionsStats struct {
	Skipped int `json:"skipped"`
}

type Relays struct {
	ApplianceType  string        `json:"appliance_type"`
	AutoOff        float64       `json:"auto_off"`
	AutoOn         float64       `json:"auto_on"`
	DefaultState   string        `json:"default_state"`
	HasTimer       bool          `json:"has_timer"`
	Ison           bool          `json:"ison"`
	MaxPower       int           `json:"max_power"`
	Name           string        `json:"name"`
	Overpower      bool          `json:"overpower"`
	Schedule       bool          `json:"schedule"`
	ScheduleRules  []interface{} `json:"schedule_rules"`
	Source         string        `json:"source"`
	TimerDuration  int           `json:"timer_duration"`
	TimerRemaining int           `json:"timer_remaining"`
	TimerStarted   int           `json:"timer_started"`
}

type Meters struct {
	Power     float64   `json:"power"`
	Overpower float64   `json:"overpower"`
	IsValid   bool      `json:"is_valid"`
	Timestamp int       `json:"timestamp"`
	Counters  []float64 `json:"counters"`
	Total     int       `json:"total"`
}
type Tmp struct {
	TC      float64 `json:"tC"`
	TF      float64 `json:"tF"`
	IsValid bool    `json:"is_valid"`
}
type Update struct {
	Status     string `json:"status"`
	HasUpdate  bool   `json:"has_update"`
	NewVersion string `json:"new_version"`
	OldVersion string `json:"old_version"`
}

type ShellySettings struct {
	Device           Device    `json:"device"`
	WifiAp           WifiAp    `json:"wifi_ap"`
	WifiSta          WifiSta   `json:"wifi_sta"`
	WifiSta1         WifiSta   `json:"wifi_sta1"`
	ApRoaming        ApRoaming `json:"ap_roaming"`
	Mqtt             Mqtt      `json:"mqtt"`
	Coiot            Coiot     `json:"coiot"`
	Sntp             Sntp      `json:"sntp"`
	Login            Login     `json:"login"`
	PinCode          string    `json:"pin_code"`
	Name             string    `json:"name"`
	Fw               string    `json:"fw"`
	Discoverable     bool      `json:"discoverable"`
	BuildInfo        BuildInfo `json:"build_info"`
	Cloud            Cloud     `json:"cloud"`
	Timezone         string    `json:"timezone"`
	Lat              float64   `json:"lat"`
	Lng              float64   `json:"lng"`
	Tzautodetect     bool      `json:"tzautodetect"`
	TzUtcOffset      int       `json:"tz_utc_offset"`
	TzDst            bool      `json:"tz_dst"`
	TzDstAuto        bool      `json:"tz_dst_auto"`
	Time             string    `json:"time"`
	Unixtime         int       `json:"unixtime"`
	LedStatusDisable bool      `json:"led_status_disable"`
	DebugEnable      bool      `json:"debug_enable"`
	AllowCrossOrigin bool      `json:"allow_cross_origin"`
	Actions          Actions   `json:"actions"`
	Hwinfo           Hwinfo    `json:"hwinfo"`
	MaxPower         int       `json:"max_power"`
	LedPowerDisable  bool      `json:"led_power_disable"`
	Relays           []Relays  `json:"relays"`
	EcoModeEnabled   bool      `json:"eco_mode_enabled"`
}

type Device struct {
	Type       string `json:"type"`
	Mac        string `json:"mac"`
	Hostname   string `json:"hostname"`
	NumOutputs int    `json:"num_outputs"`
	NumMeters  int    `json:"num_meters"`
}

type WifiAp struct {
	Enabled bool   `json:"enabled"`
	Ssid    string `json:"ssid"`
	Key     string `json:"key"`
}

type ApRoaming struct {
	Enabled   bool `json:"enabled"`
	Threshold int  `json:"threshold"`
}

type Mqtt struct {
	Enable              bool    `json:"enable"`
	Connected           bool    `json:"connected"`
	Server              string  `json:"server"`
	User                string  `json:"user"`
	ID                  string  `json:"id"`
	ReconnectTimeoutMax float64 `json:"reconnect_timeout_max"`
	ReconnectTimeoutMin float64 `json:"reconnect_timeout_min"`
	CleanSession        bool    `json:"clean_session"`
	KeepAlive           int     `json:"keep_alive"`
	MaxQos              int     `json:"max_qos"`
	Retain              bool    `json:"retain"`
	UpdatePeriod        int     `json:"update_period"`
}
type Coiot struct {
	Enabled      bool   `json:"enabled"`
	UpdatePeriod int    `json:"update_period"`
	Peer         string `json:"peer"`
}
type Sntp struct {
	Server  string `json:"server"`
	Enabled bool   `json:"enabled"`
}
type Login struct {
	Enabled     bool   `json:"enabled"`
	Unprotected bool   `json:"unprotected"`
	Username    string `json:"username"`
}
type BuildInfo struct {
	BuildID        string    `json:"build_id"`
	BuildTimestamp time.Time `json:"build_timestamp"`
	BuildVersion   string    `json:"build_version"`
}
type Actions struct {
	Active bool     `json:"active"`
	Names  []string `json:"names"`
}
type Hwinfo struct {
	HwRevision string `json:"hw_revision"`
	BatchID    int    `json:"batch_id"`
}
