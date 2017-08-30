package ipfsapp

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"os"
	"path/filepath"
)

//var dataShards = flag.Int("data", 10, "Number of shards to split the data into")
//var parShards = flag.Int("par", 4, "Number of parity shards")
var outSuffix = "pieces/"
var resSuffix = "restore/"

func fname2tmpdir(fname string) string {
	return "./tmp/" + fname + outSuffix
}

func fname2resdir(fname string) string {
	return "./" + fname + resSuffix
}
func tmpdir2fname(tmpdir string) string {
	_, file := filepath.Split(tmpdir)
	return file[:(len(file) - len(outSuffix))]
}
func rsdecode(fname string, dataShards, parShards int, cryptoKey ...string) {

	// Create matrix
	enc, err := reedsolomon.NewStream(dataShards, parShards)
	checkErr(err)
	//piecesFolder := fname2tmpdir(fname)
	// Open the inputs
	shards, size, err := openInput(dataShards, parShards, fname)
	checkErr(err)

	// Verify the shards
	ok, err := enc.Verify(shards)
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data")
		shards, size, err = openInput(dataShards, parShards, fname)
		checkErr(err)
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				outfn := fmt.Sprintf("%s.%d", fname, i)
				fmt.Println("Creating", outfn)
				out[i], err = os.Create(outfn)
				checkErr(err)
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				checkErr(err)
			}
		}
		shards, size, err = openInput(dataShards, parShards, fname)
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
			os.Exit(1)
		}
		checkErr(err)
	}

	// Join the shards and write them
	_, err = os.Open(fname2resdir(fname))
	if err != nil {
		os.Mkdir(fname2resdir(fname), os.ModeDir)
	}
	outfn := fname2resdir(fname) + fname

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	checkErr(err)

	shards, size, err = openInput(dataShards, parShards, fname)
	checkErr(err)

	// We don't know the exact filesize.
	err = enc.Join(f, shards, int64(dataShards)*size)
	checkErr(err)
}

func openInput(dataShards, parShards int, fname string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	tmpDir := fname2tmpdir(fname)
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		infn := fmt.Sprintf("%s%s.%d", tmpDir, fname, i)
		fmt.Println("Opening", infn)
		f, err := os.Open(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		checkErr(err)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}

/**
func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
**/
