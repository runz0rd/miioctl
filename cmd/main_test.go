package main

import "testing"

func Test_run(t *testing.T) {
	type args struct {
		config string
		aqi    bool
		power  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"phoney", args{"../test.yaml", true, "on"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.config, tt.args.aqi, tt.args.power); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
