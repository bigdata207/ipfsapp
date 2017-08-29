package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func ipfsInit() {
	cmd := exec.Command("ipfs", "init")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsStart() {
	cmd := exec.Command("ipfs", "daemon")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsGet(ipfsHash string) {
	cmd := exec.Command("ipfs", "get", ipfsHash)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsAddFolder(dir string) {
	cmd := exec.Command("ipfs add", "-r", dir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsAdd(fname string) {
	cmd := exec.Command("ipfs", "add", fname)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsStats(arg string) {
	args := []string{"bitswap", "bw", "repo"}
	b := containEle(strArr2InterArr(args), arg)
	if b == -1 {
		return
	}
	cmd := exec.Command("ipfs", "stat", arg)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}
