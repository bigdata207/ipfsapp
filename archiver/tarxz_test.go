package archiver

import "testing"

func Test_xzFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		x    xzFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := xzFormat{}
			if got := x.Match(tt.args.filename); got != tt.want {
				t.Errorf("xzFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTarXz(t *testing.T) {
	type args struct {
		tarxzPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTarXz(tt.args.tarxzPath); got != tt.want {
				t.Errorf("isTarXz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_xzFormat_Make(t *testing.T) {
	type args struct {
		xzPath    string
		filePaths []string
	}
	tests := []struct {
		name    string
		x       xzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := xzFormat{}
			if err := x.Make(tt.args.xzPath, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("xzFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_xzFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		x       xzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := xzFormat{}
			if err := x.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("xzFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
