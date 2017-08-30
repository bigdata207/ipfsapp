package ipfsapp

import (
	"fmt"
	"testing"
)

func TestLinesCounter(t *testing.T) {
	fmt.Println("Test LinesConter()")
	p := "/mnt/extra/MackZ/go/src/github.com/mackzhong/ipfsapp"
	fmt.Printf("lines of code in %v : %d\n", p, LinesCounter(p))
	registerFunc("strend", strend)
	fun := func2anonymousFunc("strend")
	a := "i am student!"
	fmt.Println(string(fun(a).(byte)))
}
