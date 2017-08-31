package archiver

import "testing"

func Test_sevenzFormat_Match(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		sz   sevenzFormat
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sz := sevenzFormat{}
			if got := sz.Match(tt.args.filename); got != tt.want {
				t.Errorf("sevenzFormat.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sevenzFormat_Make(t *testing.T) {
	type args struct {
		destination string
		sources     []string
	}
	tests := []struct {
		name    string
		sz      sevenzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sz := sevenzFormat{}
			if err := sz.Make(tt.args.destination, tt.args.sources); (err != nil) != tt.wantErr {
				t.Errorf("sevenzFormat.Make() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sevenzFormat_Open(t *testing.T) {
	type args struct {
		source      string
		destination string
	}
	tests := []struct {
		name    string
		sz      sevenzFormat
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sz := sevenzFormat{}
			if err := sz.Open(tt.args.source, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("sevenzFormat.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isSevenZ(t *testing.T) {
	type args struct {
		sevenzPath string
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
			if got := isSevenZ(tt.args.sevenzPath); got != tt.want {
				t.Errorf("isSevenZ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sevenzFile(t *testing.T) {
	type args struct {
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
			if err := sevenzFile(tt.args.source); (err != nil) != tt.wantErr {
				t.Errorf("sevenzFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_unsevenzFile(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := unsevenzFile(); (err != nil) != tt.wantErr {
				t.Errorf("unsevenzFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
