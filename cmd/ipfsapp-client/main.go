package main

import (
	"flag"
	"github.com/mackzhong/ipfsapp/ipfscmd"
	"os"
)

func main() {
	args := flag.Args()
	Args := append([]string{os.Args[0]}, args...)
	ret := ipfscmd.MainRet(Args...)
	if ret
}
