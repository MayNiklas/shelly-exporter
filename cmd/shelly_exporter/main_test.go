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

func makeMockShelly(settings, status string) *httptest.Server {

	// Create fake server (fake shelly endpoint) that responds json from our
	// testdata files
	mux := http.NewServeMux()
	mux.HandleFunc("/settings", makeMockHandler(settings))
	mux.HandleFunc("/status", makeMockHandler(status))
	return httptest.NewServer(mux)

}

func Test_probeHandler(t *testing.T) {

	tests := []struct {
		name             string
		expectedFile     string
		settingsJsonFile string
		statusJsonFile   string
	}{
		{
			name:             "first test",
			settingsJsonFile: "../../tests/settings.json",
			statusJsonFile:   "../../tests/status.json",
			expectedFile:     "../../tests/metrics.prom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockShellySrv := makeMockShelly(tt.settingsJsonFile, tt.statusJsonFile)
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

			content, err := ioutil.ReadFile(tt.expectedFile)
			if err != nil {
				panic(err)
			}
			metrics := string(content)

			if string(data) != metrics {
				t.Errorf("Expected %v, got: \n%v", metrics, string(data))
			}
		})
	}
}
