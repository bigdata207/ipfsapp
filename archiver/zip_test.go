package archiver

import (
	"archive/zip"
	"testing"
)

func Test_zipFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		z    zipFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := zipFormat{}
			if got := z.Match(tt.args.filename); got != tt.want {
				t.Errorf("zipFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isZip(t *testing.T) {
	type args struct {
		zipPath string
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
			if got := isZip(tt.args.zipPath); got != tt.want {
				t.Errorf("isZip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zipFormat_Make(t *testing.T) {
	type args struct {
		zipPath   string
		filePaths []string
	}
	tests := []struct {
		name    string
		z       zipFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := zipFormat{}
			if err := z.Make(tt.args.zipPath, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("zipFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_zipFile(t *testing.T) {
	type args struct {
		w      *zip.Writer
		source string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := zipFile(tt.args.w, tt.args.source); (err != nil) != tt.wantErr {
				t.Errorf("zipFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_zipFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		z       zipFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := zipFormat{}
			if err := z.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("zipFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_unzipFile(t *testing.T) {
	type args struct {
		zf          *zip.File
		destination string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := unzipFile(tt.args.zf, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("unzipFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
