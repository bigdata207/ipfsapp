package main

import (
	"github.com/mackzhong/ipfsapp/ipfscmd"
)

func startIPFS(c chan int) error {
	Args := append([]string{os.Args[0]}, "daemon")
	go ipfscmd.MainRet(c, Args...)
}
