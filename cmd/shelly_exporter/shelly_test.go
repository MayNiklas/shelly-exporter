package shelly_exporter

import (
	"reflect"
	"testing"
)

func Test_getJson(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getJson(tt.args.url)
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
	type fields struct {
		shelly_status   shelly_status
		shelly_settings shelly_settings
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ShellyData{
				shelly_status:   tt.fields.shelly_status,
				shelly_settings: tt.fields.shelly_settings,
			}
			if err := s.Fetch(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("ShellyData.Fetch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
