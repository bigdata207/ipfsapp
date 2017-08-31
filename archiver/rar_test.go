package archiver

import "testing"

func Test_rarFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		r    rarFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rarFormat{}
			if got := r.Match(tt.args.filename); got != tt.want {
				t.Errorf("rarFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRar(t *testing.T) {
	type args struct {
		rarPath string
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
			if got := isRar(tt.args.rarPath); got != tt.want {
				t.Errorf("isRar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rarFormat_Make(t *testing.T) {
	type args struct {
		rarPath   string
		filePaths []string
	}
	tests := []struct {
		name    string
		r       rarFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rarFormat{}
			if err := r.Make(tt.args.rarPath, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("rarFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_rarFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		r       rarFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rarFormat{}
			if err := r.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("rarFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
