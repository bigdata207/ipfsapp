package main

import (
	"fmt"
	"reflect"
)

var structMap map[string]reflect.Type

func registerStruct(obj interface{}) {
	rv := reflect.TypeOf(obj)
	fmt.Println(rv.Name())
	structMap[rv.Name()] = rv
	fmt.Println(structMap[rv.Name()].NumField())
}

func registerInitStruct() {
	structMap = make(map[string]reflect.Type)
}
