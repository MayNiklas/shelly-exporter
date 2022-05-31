package shelly_exporter

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"io/ioutil"
	"reflect"
	"testing"
)

func Test_getJson(t *testing.T) {

	mockShelly := makeMockShelly("../../tests/settings.json", "../../tests/status.json")

	settingsJsonBytes, err := ioutil.ReadFile("../../tests/settings.json")
	if err != nil {
		panic(err)
	}
	// settingsJson := string(content)

	statusJsonBytes, err := ioutil.ReadFile("../../tests/status.json")
	if err != nil {
		panic(err)
	}
	// statusJson := string(content)

	tests := []struct {
		name    string
		url     string
		want    []byte
		wantErr bool
	}{
		{
			name:    "Test fail connection",
			url:     "http://nowhere.json",
			want:    nil,
			wantErr: true,
		},
		{
			name: "GET status.json",
			url:  mockShelly.URL + "/status",
			want: statusJsonBytes,
			wantErr: false,
		},
		{
			name: "GET settings.json",
			url:  mockShelly.URL + "/settings",
			want: settingsJsonBytes,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJson(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJson(%v) error = %v, wantErr %v", tt.url, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getJson(%v) = %v, want %v",tt.url, string(got), string(tt.want))
			}
		})
	}
}

func TestShellyData_Fetch(t *testing.T) {

	mockShelly := makeMockShelly("../../tests/settings.json", "../../tests/status.json")

	exampleShellyData := ShellyData{
		Settings: ShellySettings{
			Device: Device{
				Type:       "SHPLG-S",
				Mac:        "8CAAB9EDE2EE",
				Hostname:   "shellyplug-s-EAE4EE",
				NumOutputs: 1,
				NumMeters:  1,
			},
			WifiAp:    WifiAp{Ssid: "shellyplug-s-EAE4EE"},
			WifiSta:   WifiSta{Enabled: true, Ssid: "i-wont-tell", Ipv4Method: "dhcp"},
			WifiSta1:  WifiSta{Ipv4Method: "dhcp"},
			ApRoaming: ApRoaming{Threshold: -70},
			Mqtt: Mqtt{
				Server:              "192.168.33.3:1883",
				ID:                  "shellyplug-s-EAE4EE",
				ReconnectTimeoutMax: 60,
				ReconnectTimeoutMin: 2,
			},
			Coiot: Coiot{Enabled: true, UpdatePeriod: 15},
			Login: Login{
				Enabled:     false,
				Unprotected: false,
				Username:    "admin",
			},
			Fw: "20220209-094058/v1.11.8-g8c7bb8d",
			Hwinfo: Hwinfo{
				HwRevision: "prod-190516",
				BatchID:    1,
			},

			Actions: Actions{
				Active: false,
				Names:  []string{"btn_on_url", "out_on_url", "out_off_url"},
			},
			EcoModeEnabled: true,
			MaxPower:       2500,
			PinCode:        "w<#LkL",
			Name:           "#1 Rack",
			Sntp:           Sntp{Server: "time.google.com", Enabled: true},
			Timezone:       "Europe/Berlin",
			Lat:            51.31477,
			Lng:            9.49086,
			Tzautodetect:   true,
			TzUtcOffset:    7200,
			TzDstAuto:      true,
			Time:           "16:30",
			Unixtime:       1653921059,
		},

		Status: ShellyStatus{
			RAMTotal:      51264,
			RAMFree:       38632,
			FsSize:        233681,
			FsFree:        165158,
			Uptime:        7767857,
			Time:          "16:30",
			Temperature:   31.44,
			WifiSta:       WifiSta{Ssid: "i-wont-tell", Connected: true, IP: "192.168.15.2", Rssi: -67},
			Cloud:         Cloud{Enabled: true, Connected: true},
			Unixtime:      1653921059,
			Serial:        31519,
			Mac:           "8CAAB9EDE2EE",
			CfgChangedCnt: 4,
			Relays:        []Relays{{}},
		},
	}

	tests := []struct {
		name    string
		address string
		want    ShellyData
		wantErr bool
	}{
		{
			name:    "Test fetch fail",
			address: "127.0.0.1",
			want:    exampleShellyData,
			wantErr: true,
		},
		{
			name:    "Test fetch success",
			want:    exampleShellyData,
			address: mockShelly.URL,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := ShellyData{}

			var err error

			if err = s.Fetch(tt.address); (err != nil) != tt.wantErr {
				t.Errorf("ShellyData.Fetch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {

				// Ignore some values so we don't have to specify everything in
				// the tests
				opts := []cmp.Option{
					cmpopts.IgnoreFields(ShellyStatus{}, "Relays", "Meters", "Tmp", "Update"),
					cmpopts.IgnoreFields(ShellySettings{}, "Mqtt", "Relays", "BuildInfo", "Cloud"),
				}

				if diff := cmp.Diff(tt.want, s, opts...); diff != "" {
					t.Errorf("ShellyData (after Fetch) mismatch (-want +got):\n%s", diff)
				}
			}

		})
	}
}
