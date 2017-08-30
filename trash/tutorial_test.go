package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type USB interface {
	Name() string
	Connect()
}

type PhoneConnector struct {
	name string
}

func (p PhoneConnector) Name() string {
	return p.name
}

func (p PhoneConnector) Connect() {
	fmt.Println("Connect:", p.name)
}

type person struct {
	Name string
	Age  int
}

type student struct {
	person
	StuID string
}

func (person) sayHi() {
	fmt.Println("Hi")
}

func (p *person) setAge() {
	func(p *person) {
		p.Age = 18
	}(p)
}

func getInfo(o interface{}) {
	t := reflect.TypeOf(o)
	fmt.Println("Type:", t.Name())

	v := reflect.ValueOf(o)
	fmt.Println("Fields: ")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Printf("%6s: %v = %v\n", f.Name, f.Type, val)
	}
}

type Citizen struct {
	ID   int
	Name string
	Age  int
}

func (Citizen) Hello() {
	fmt.Println("Hello World!")
}

func mapTest() {
	fmt.Println("hello")
	a := make(map[string]string)
	a["me"] = "zt"
	for k, v := range a {
		fmt.Printf("%s: %6s\n", k, v)
	}
	b := map[int]string{1: "zt", 2: "sustc"}
	for k, v := range b {
		fmt.Printf("%v: %6s\n", k, v)
	}
}
func arrayTest() {
	var a = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	var b = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	var j int
	for j = 0; j < len(a); j++ {
		fmt.Printf("a[%d] = %v\n", j, a[j])
	}
	for j = 0; j < len(b); j++ {
		fmt.Printf("b[%d] = %v\n", j, b[j])
	}
	c := [3][4]int{
		{0, 1, 2, 3},   /*  第一行索引为 0 */
		{4, 5, 6, 7},   /*  第二行索引为 1 */
		{8, 9, 10, 11}, /*  第三行索引为 2 */
	}
	for _, col := range c {
		for _, val := range col {
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
	fmt.Println(c[1:3])
	d := make([]int, 0, 100)
	d = append(d, 12, 13, 14, 15)
	fmt.Println(d)
	e := make([]int, len(d), cap(d))
	copy(e, d)
	fmt.Println(e)
}
func controlTest() {
	if true {
		fmt.Println("Test if")
		if false {

		} else {
			fmt.Println("Test else")
		}
	}

	a := 2
	fmt.Print("Test switch a = ")
	switch a {
	case 0:
		fmt.Println(a)
		break
	case 1:
		fmt.Println(a)
		break
	case 2:
		fmt.Println(a)
		break
	default:
		fmt.Println(a)
		break
	}
	var c1, c2, c3 chan int
	var i1, i2 int
	select {
	case i1 = <-c1:
		fmt.Printf("received ", i1, " from c1\n")
	case c2 <- i2:
		fmt.Printf("sent ", i2, " to c2\n")
	case i3, ok := (<-c3): // same as: i3, ok := <-c3
		if ok {
			fmt.Printf("received ", i3, " from c3\n")
		} else {
			fmt.Printf("c3 is closed\n")
		}
	default:
		fmt.Printf("no communication\n")
	}
}

func pointTest() {
	var a int = 20 /* 声明实际变量 */
	var ip *int    /* 声明指针变量 */

	ip = &a /* 指针变量的存储地址 */

	fmt.Printf("a 变量的地址是: %x\n", &a)

	/* 指针变量的存储地址 */
	fmt.Printf("ip 变量储存的指针地址: %x\n", ip)

	/* 使用指针访问值 */
	fmt.Printf("*ip 变量的值: %d\n", *ip)
}

func TestTutorial(t *testing.T) {
	mapTest()
	arrayTest()
	controlTest()
	pointTest()
	p := &person{"zt", 20}
	a := struct {
		Name string
		Age  int
	}{
		Name: "sustc",
		Age:  6,
	}
	fmt.Println(a)
	p.sayHi()
	fmt.Println(p)
	p.setAge()
	fmt.Println(p)
	stu := &student{person: person{"zt", 17}, StuID: "11510291"}
	stu.sayHi()
	fmt.Println(stu)
	stu.Name = "sustc"
	fmt.Println(stu)
	time.Sleep(time.Second)
	for i := 0; i < 3; i++ {
		fmt.Println("i =", i)
		defer fmt.Println("i =", i)
		defer func() {
			fmt.Println("Niming: ", i)
		}()
	}
	type id int
	k := id(1)
	fmt.Println(k)
	c := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Go")
			c <- i
		}
		close(c)
	}()
	for u := range c {
		fmt.Printf("sub: %v\n", u)
	}
	var b USB
	b = PhoneConnector{name: "MeiZu"}
	b.Connect()
	me := Citizen{1, "zt", 12}
	go getInfo(me)
	time.Sleep(time.Second)
	fmt.Println("Final:")
}
