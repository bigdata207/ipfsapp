package ipfsapp

import (
	"fmt"
	"testing"
)

func TestWebsite(t *testing.T) {
	fmt.Println("Test Run Vue2.0 Website")
	c := make(chan error)
	npmInstall(c)
	if err := <-c; err != nil {
		fmt.Println(err.Error())
		return
	}
	npmStart(c)
	if <-c == nil {
		npmRunDev(c)
		<-c
	}
}
