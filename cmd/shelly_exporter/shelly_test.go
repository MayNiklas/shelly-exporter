package shelly_exporter

import (
	"reflect"
	"testing"
)

func Test_getJson(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    []byte
		wantErr bool
	}{
		{
			name: "Test fail",
			url:  "TODO",
			want: []byte{
				//TODO
			},
			wantErr: true,
		},
		{
			name: "Test success",
			url:  "TODO",
			want: []byte{
				//TODO
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJson(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShellyData_Fetch(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		Status   ShellyStatus
		Settings ShellySettings
		wantErr  bool
	}{
		{
			name:     "Test fail",
			address:  "127.0.0.1",
			Settings: ShellySettings{
				//TODO
			},
			Status: ShellyStatus{
				//TODO
			},
			wantErr: true,
		},
		{
			name:     "Test success",
			address:  "127.0.0.1",
			Settings: ShellySettings{},
			Status:   ShellyStatus{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShellyData{
				Status:   tt.Status,
				Settings: tt.Settings,
			}
			if err := s.Fetch(tt.address); (err != nil) != tt.wantErr {
				t.Errorf("ShellyData.Fetch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
