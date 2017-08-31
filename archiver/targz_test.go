package archiver

import "testing"

func Test_tarGzFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		t    tarGzFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarGzFormat{}
			if got := t.Match(tt.args.filename); got != tt.want {
				t.Errorf("tarGzFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTarGz(t *testing.T) {
	type args struct {
		targzPath string
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
			if got := isTarGz(tt.args.targzPath); got != tt.want {
				t.Errorf("isTarGz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarGzFormat_Make(t *testing.T) {
	type args struct {
		targzPath string
		filePaths []string
	}
	tests := []struct {
		name    string
		t       tarGzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarGzFormat{}
			if err := t.Make(tt.args.targzPath, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("tarGzFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarGzFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		t       tarGzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarGzFormat{}
			if err := t.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("tarGzFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
