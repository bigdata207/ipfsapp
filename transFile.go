package ipfsapp

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var bs int = 1024 * 512

type chanPair struct {
	in  chan string
	out chan string
}

func NewChanPair() *chanPair {
	cp := &chanPair{}
	cp.in = make(chan string)
	cp.out = make(chan string)
	return cp
}

func goTransFile(origFileName string) error {
	t1 := time.Now()
	gorun := make(map[string]*chanPair)
	a := "A"
	//origFileName := "xp.iso"
	gorun[origFileName] = NewChanPair()
	buf := make([]byte, bs)
	go func(cp *chanPair) {
		fileName := <-cp.in
		fmt.Println("Got FileName: ", fileName)
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Out file create failed. err: " + err.Error())
		}
		defer f.Close()
		end := 0
	loop:
		for {
			select {
			case v := <-cp.in:
				if strings.Compare("finish", v) == 0 {
					fmt.Println(v + ":" + a)
					cp.out <- "All Right"
					break loop
				} else {
					l, err := strconv.Atoi(v)
					if err == nil && l > 0 {
						//err = appendToFile(fileName, buf[:l])
						k, err := f.WriteAt(buf[:l], int64(end))
						end += k
						if err != nil || k != l {
							fmt.Println("Write Error!")
							cp.out <- "finish"
						}
						cp.out <- "ok"
					} else {
						fmt.Println("Don't know the buf size")
					}

				}
			}
		}
		fmt.Println("End")
	}(gorun[origFileName])

	orig, err := os.OpenFile(origFileName, os.O_RDONLY, 0644)
	gorun[origFileName].in <- "out." + strings.Split(origFileName, ".")[1]
	if err != nil {
		fmt.Println("Open File Error")
		return errors.New("Open File Error")
	}
	for kk, err := orig.Read(buf); err == nil; kk, err = orig.Read(buf) {
		//fmt.Println("Read size: ", kk)
		gorun[origFileName].in <- strconv.Itoa(kk)

		res := <-gorun[origFileName].out
		//fmt.Println("Wait: ", res)
		if strings.Compare(res, "finish") == 0 {
			break
		}
	}
	gorun[origFileName].in <- "finish"
	//time.Sleep(time.Millisecond)
	fmt.Println("Final: ", <-gorun[origFileName].out)
	time.Sleep(time.Millisecond * 10)
	fmt.Println(time.Since(t1))
	return nil
}
