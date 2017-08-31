package archiver

import "testing"

func Test_tarBz2Format_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		t    tarBz2Format
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarBz2Format{}
			if got := t.Match(tt.args.filename); got != tt.want {
				t.Errorf("tarBz2Format.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTarBz2(t *testing.T) {
	type args struct {
		tarbz2Path string
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
			if got := isTarBz2(tt.args.tarbz2Path); got != tt.want {
				t.Errorf("isTarBz2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarBz2Format_Make(t *testing.T) {
	type args struct {
		tarbz2Path string
		filePaths  []string
	}
	tests := []struct {
		name    string
		t       tarBz2Format
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarBz2Format{}
			if err := t.Make(tt.args.tarbz2Path, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("tarBz2Format.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarBz2Format_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		t       tarBz2Format
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarBz2Format{}
			if err := t.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("tarBz2Format.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
