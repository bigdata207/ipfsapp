package archiver

import "testing"

func Test_tarLz4Format_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		t    tarLz4Format
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarLz4Format{}
			if got := t.Match(tt.args.filename); got != tt.want {
				t.Errorf("tarLz4Format.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTarLz4(t *testing.T) {
	type args struct {
		tarlz4Path string
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
			if got := isTarLz4(tt.args.tarlz4Path); got != tt.want {
				t.Errorf("isTarLz4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarLz4Format_Make(t *testing.T) {
	type args struct {
		tarlz4Path string
		filePaths  []string
	}
	tests := []struct {
		name    string
		t       tarLz4Format
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarLz4Format{}
			if err := t.Make(tt.args.tarlz4Path, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("tarLz4Format.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarLz4Format_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		t       tarLz4Format
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarLz4Format{}
			if err := t.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("tarLz4Format.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
