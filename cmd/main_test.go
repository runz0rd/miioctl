package main

import "testing"

func Test_run(t *testing.T) {
	type args struct {
		config    string
		status    string
		serveAddr string
		power     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// {"status", args{config: "../example.yaml", status: "pm25"}, false},
		// {"power", args{config: "../example.yaml", power: "off"}, false},
		{"server", args{config: "../example.yaml", serveAddr: ":8080"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := run(tt.args.config, tt.args.status, tt.args.serveAddr, tt.args.power); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
