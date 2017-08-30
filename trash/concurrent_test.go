package mackzhong

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func Test_Concurrent(t *testing.T) {
	fmt.Println("Test Concurrent")
	defer fmt.Println("结束啦")
	//闭包的话defer的指令里包含的外部变量将会使用最新的值，而不是调用闭包时的值,
	//但是取外部变量值的时候时取当时的值
	for i := 0; i < 4; i++ {
		closure := func() func() {
			j := i
			return func() {
				fmt.Printf("i,j = %d,%d\n", i, j)
			}
		}()
		defer closure()
	}

	//协程通信
	chs := make([]chan int, 5)
	//需要初始化chs数组
	for i := range chs {
		fmt.Println(i)
		chs[i] = make(chan int)
	}
	go func() {
	Loop:
		for {
			select {
			case v1 := <-chs[0]:
				fmt.Printf("Get %v from chs[0]\n", v1)
			case v2 := <-chs[1]:
				fmt.Printf("Get %v from chs[1]\n", v2)
			case v3 := <-chs[2]:
				fmt.Printf("Get %v from chs[2]\n", v3)
			case <-chs[3]:
				break Loop
			}
			fmt.Println("啥情况了")
		}
		chs[4] <- 1
	}()
	for i := 0; i < 4; i++ {
		chs[i] <- i
	}

	<-chs[4]
	close(chs[0])
	close(chs[1])
	close(chs[2])
	close(chs[3])
	close(chs[4])

	//遍历chan数组
	GetChan := func(ch chan int16) {
		ch <- 10086
	}
	chss := make([]chan int16, 5)
	for i := 0; i < 5; i++ {
		chss[i] = make(chan int16)
		go GetChan(chss[i])
	}

	//缓冲
	k := make(chan int, 8)
	for i := 0; i < 4; i++ {
		k <- i
	}
	//需要关闭输入防止死锁
	close(k)
	for v := range k {
		fmt.Printf("缓冲区内容: %d\n", v)
	}
	for _, ch := range chss {
		fmt.Println(<-ch)
	}
	//单向channel
	c := make(chan int, 3)

	//不知道为啥:= chan<- int(c)这种方式不行....
	var send chan<- int = c
	var recv <-chan int = c
	//关闭c时已经关闭send,recv是只读channel，不用关闭
	defer close(c)
	send <- 123

	println(<-recv)
	//锁测试
	var counter int
	counter = 0
	lock := &sync.Mutex{}
	Count := func(lock *sync.Mutex) {
		lock.Lock()
		counter++
		fmt.Println(counter)
		lock.Unlock()
	}
	for i := 0; i < 5; i++ {
		go Count(lock)
	}
	for {
		lock.Lock()
		c := counter
		lock.Unlock()
		runtime.Gosched()
		if c >= 5 {
			break
		}
	}

	//多核并行
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())

	//全局唯一性操作
	var a string
	var once sync.Once
	setup := func() {
		fmt.Println("Setup a")
		a = "hello go"
	}

	//里面的只会执行一次setup,当其他协程发现已经有进程运行setup时会阻塞直
	//到setup调用结束，之后的协程不会再调用setup
	doprint := func() {
		once.Do(setup)
		fmt.Println(a)
	}
	go doprint()
	go doprint()
	//存在通信延迟
	time.Sleep(time.Microsecond * 10)
	fmt.Println("a = ", a)
	//嵌套协程
	for i := 0; i < 3; i++ {
		go func() {
			fmt.Println("Go")
			go func() {
				fmt.Println("GoGo")
			}()
		}()
	}
}
