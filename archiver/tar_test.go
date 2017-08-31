package archiver

import (
	"archive/tar"
	"testing"
)

func Test_tarFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		t    tarFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarFormat{}
			if got := t.Match(tt.args.filename); got != tt.want {
				t.Errorf("tarFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isTar(t *testing.T) {
	type args struct {
		tarPath string
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
			if got := isTar(tt.args.tarPath); got != tt.want {
				t.Errorf("isTar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasTarHeader(t *testing.T) {
	type args struct {
		buf []byte
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
			if got := hasTarHeader(tt.args.buf); got != tt.want {
				t.Errorf("hasTarHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tarFormat_Make(t *testing.T) {
	type args struct {
		tarPath   string
		filePaths []string
	}
	tests := []struct {
		name    string
		t       tarFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarFormat{}
			if err := t.Make(tt.args.tarPath, tt.args.filePaths); (err != nil) != tt.wantErr {
				t.Errorf("tarFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarball(t *testing.T) {
	type args struct {
		filePaths []string
		tarWriter *tar.Writer
		dest      string
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
			if err := tarball(tt.args.filePaths, tt.args.tarWriter, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("tarball() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarFile(t *testing.T) {
	type args struct {
		tarWriter *tar.Writer
		source    string
		dest      string
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
			if err := tarFile(tt.args.tarWriter, tt.args.source, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("tarFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tarFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		t       tarFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := tarFormat{}
			if err := t.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("tarFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_untar(t *testing.T) {
	type args struct {
		tr          *tar.Reader
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
			if err := untar(tt.args.tr, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("untar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_untarFile(t *testing.T) {
	type args struct {
		tr          *tar.Reader
		header      *tar.Header
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
			if err := untarFile(tt.args.tr, tt.args.header, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("untarFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
