package ipfsapp

import (
	"fmt"
	"os"
	_ "strings"
	"testing"
	"time"
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

func TestGetIP(t *testing.T) {
	getInternalIP()
}

func TestIPAPI(t *testing.T) {
	t1 := time.Now()
	fmt.Println(IPAPI(true))
	fmt.Println(time.Since(t1))
	fmt.Println(os.Getwd())
}
