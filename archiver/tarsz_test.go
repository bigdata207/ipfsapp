package archiver

import "testing"

func Test_tarSzFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		t    tarSzFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarSzFormat{}
			if got := t.Match(tt.args.filename); got != tt.want {
				t.Errorf("tarSzFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTarSz(t *testing.T) {
	type args struct {
		tarszPath string
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
			if got := isTarSz(tt.args.tarszPath); got != tt.want {
				t.Errorf("isTarSz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarSzFormat_Make(t *testing.T) {
	type args struct {
		tarszPath string
		filePaths []string
	}
	tests := []struct {
		name    string
		t       tarSzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarSzFormat{}
			if err := t.Make(tt.args.tarszPath, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("tarSzFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarSzFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		t       tarSzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarSzFormat{}
			if err := t.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("tarSzFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
