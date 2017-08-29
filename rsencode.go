package main

import (
	"flag"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
	"os"
	"path/filepath"
)

var dataShards = flag.Int("data", 10, "Number of shards to split the data into, must be below 257.")
var parShards = flag.Int("par", 4, "Number of parity shards")
var outDir = flag.String("out", "/mnt/extra/MackZ/go/src/github.com/mackzhong/ec/output", "Alternative output directory")

func rsencode() {
	fmt.Println("Test ReedSolomon")
	// Parse command line parameters.
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		fmt.Fprintf(os.Stderr, "Error: No input filename given\n")
		flag.Usage()
		os.Exit(1)
	}
	if *dataShards > 257 {
		fmt.Fprintf(os.Stderr, "Error: Too many data shards\n")
		os.Exit(1)
	}
	fname := "/mnt/extra/MackZ/go/src/github.com/mackzhong/ec/1.mp3"

	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(*dataShards, *parShards)
	checkErr(err)

	fmt.Println("Opening", fname)
	f, err := os.Open(fname)
	checkErr(err)

	instat, err := f.Stat()
	checkErr(err)

	shards := *dataShards + *parShards
	out := make([]*os.File, shards)

	// Create the resulting files.
	dir, file := filepath.Split(fname)
	if *outDir != "" {
		dir = *outDir
	}
	for i := range out {
		outfn := fmt.Sprintf("%s.%d", file, i)
		fmt.Println("Creating", outfn)
		out[i], err = os.Create(filepath.Join(dir, outfn))
		checkErr(err)
	}

	// Split into files.
	data := make([]io.Writer, *dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, instat.Size())
	checkErr(err)

	// Close and re-open the files.
	input := make([]io.Reader, *dataShards)

	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		checkErr(err)
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, *parShards)
	for i := range parity {
		parity[i] = out[*dataShards+i]
		defer out[*dataShards+i].Close()
	}

	// Encode parity
	err = enc.Encode(input, parity)
	checkErr(err)
	fmt.Printf("File split into %d data + %d parity shards.\n", *dataShards, *parShards)

}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(2)
	}
}
