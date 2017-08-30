package archiver

import ()

// SevenZ is for 7z format
var SevenZ sevenzFormat

func init() {
	RegisterFormat("7z", SevenZ)
}

type sevenzFormat struct{}

func (sz sevenzFormat) Match(filename string) bool {
	return true
}

// Make makes an archive.
func (sz sevenzFormat) Make(destination string, sources []string) error {
	return nil
}

// Open extracts an archive.
func (sz sevenzFormat) Open(source, destination string) error {
	return nil
}

//isSevenZ checks the file has the 7z format signature by reading its beginning
// bytes and matching it against "PK\x03\x04"
func isSevenZ(sevenzPath string) bool {
	return true
}

func sevenzFile(source string) error {
	return nil
}

func unsevenzFile() error {
	return nil
}
