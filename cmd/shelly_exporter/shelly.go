package shelly_exporter

import "time"

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

type shelly_settings struct {
	Device struct {
		Type       string `json:"type"`
		Mac        string `json:"mac"`
		Hostname   string `json:"hostname"`
		NumOutputs int    `json:"num_outputs"`
		NumMeters  int    `json:"num_meters"`
	} `json:"device"`
	WifiAp struct {
		Enabled bool   `json:"enabled"`
		Ssid    string `json:"ssid"`
		Key     string `json:"key"`
	} `json:"wifi_ap"`
	WifiSta struct {
		Enabled    bool        `json:"enabled"`
		Ssid       string      `json:"ssid"`
		Ipv4Method string      `json:"ipv4_method"`
		IP         interface{} `json:"ip"`
		Gw         interface{} `json:"gw"`
		Mask       interface{} `json:"mask"`
		DNS        interface{} `json:"dns"`
	} `json:"wifi_sta"`
	WifiSta1 struct {
		Enabled    bool        `json:"enabled"`
		Ssid       interface{} `json:"ssid"`
		Ipv4Method string      `json:"ipv4_method"`
		IP         interface{} `json:"ip"`
		Gw         interface{} `json:"gw"`
		Mask       interface{} `json:"mask"`
		DNS        interface{} `json:"dns"`
	} `json:"wifi_sta1"`
	ApRoaming struct {
		Enabled   bool `json:"enabled"`
		Threshold int  `json:"threshold"`
	} `json:"ap_roaming"`
	Mqtt struct {
		Enable              bool    `json:"enable"`
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
	} `json:"mqtt"`
	Coiot struct {
		Enabled      bool   `json:"enabled"`
		UpdatePeriod int    `json:"update_period"`
		Peer         string `json:"peer"`
	} `json:"coiot"`
	Sntp struct {
		Server  string `json:"server"`
		Enabled bool   `json:"enabled"`
	} `json:"sntp"`
	Login struct {
		Enabled     bool   `json:"enabled"`
		Unprotected bool   `json:"unprotected"`
		Username    string `json:"username"`
	} `json:"login"`
	PinCode      string `json:"pin_code"`
	Name         string `json:"name"`
	Fw           string `json:"fw"`
	Discoverable bool   `json:"discoverable"`
	BuildInfo    struct {
		BuildID        string    `json:"build_id"`
		BuildTimestamp time.Time `json:"build_timestamp"`
		BuildVersion   string    `json:"build_version"`
	} `json:"build_info"`
	Cloud struct {
		Enabled   bool `json:"enabled"`
		Connected bool `json:"connected"`
	} `json:"cloud"`
	Timezone         string  `json:"timezone"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	Tzautodetect     bool    `json:"tzautodetect"`
	TzUtcOffset      int     `json:"tz_utc_offset"`
	TzDst            bool    `json:"tz_dst"`
	TzDstAuto        bool    `json:"tz_dst_auto"`
	Time             string  `json:"time"`
	Unixtime         int     `json:"unixtime"`
	LedStatusDisable bool    `json:"led_status_disable"`
	DebugEnable      bool    `json:"debug_enable"`
	AllowCrossOrigin bool    `json:"allow_cross_origin"`
	Actions          struct {
		Active bool     `json:"active"`
		Names  []string `json:"names"`
	} `json:"actions"`
	Hwinfo struct {
		HwRevision string `json:"hw_revision"`
		BatchID    int    `json:"batch_id"`
	} `json:"hwinfo"`
	MaxPower        int  `json:"max_power"`
	LedPowerDisable bool `json:"led_power_disable"`
	Relays          []struct {
		Name          interface{}   `json:"name"`
		ApplianceType string        `json:"appliance_type"`
		Ison          bool          `json:"ison"`
		HasTimer      bool          `json:"has_timer"`
		DefaultState  string        `json:"default_state"`
		AutoOn        float64       `json:"auto_on"`
		AutoOff       float64       `json:"auto_off"`
		Schedule      bool          `json:"schedule"`
		ScheduleRules []interface{} `json:"schedule_rules"`
		MaxPower      int           `json:"max_power"`
	} `json:"relays"`
	EcoModeEnabled bool `json:"eco_mode_enabled"`
}
