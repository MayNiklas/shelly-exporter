package shelly_exporter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeMockHandler(filepath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonString, err := ioutil.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, string(jsonString))
	}
}

func Test_probeHandler(t *testing.T) {

	tests := []struct {
		name             string
		expected         string
		settingsJsonFile string
		statusJsonFile   string
	}{
		{
			name:             "first test",
			settingsJsonFile: "../../tests/settings.json",
			statusJsonFile:   "../../tests/status.json",
			expected:         "ABC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create fake server (fake shelly endpoint) that responds json from our
			// testdata files
			mux := http.NewServeMux()
			mux.HandleFunc("/settings", makeMockHandler(tt.settingsJsonFile))
			mux.HandleFunc("/status", makeMockHandler(tt.statusJsonFile))
			mockShellySrv := httptest.NewServer(mux)
			defer mockShellySrv.Close()

			// Make request to mock server on loopback interface
			request := httptest.NewRequest(http.MethodGet, "/metrics?target="+mockShellySrv.URL, nil)
			w := httptest.NewRecorder()

			probeHandler(w, request)

			res := w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)

			if err != nil {
				t.Errorf("Expected error to be nil, got %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, string(data))
			}
		})
	}
}
