package ipfsapp

import (
	"fmt"
	"testing"
)

func getFunc() func() {
	i := 0
	return func() {
		fmt.Println(i)
		i++
	}
}

func TestServer(*testing.T) {
	fmt.Println("Start RESTful API Server...")
	c := make(chan error)
	go StartAPIServer(c)

	f := getFunc()
	f()
	f()
	f()
	fmt.Println(<-c)
}
