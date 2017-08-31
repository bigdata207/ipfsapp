package ipfsapp

import (
	"fmt"
	"reflect"
)

var structMap map[string]reflect.Type

//funcMap 通过函数名获取对应的匿名函数，将API对外接口封装成基本HTTP操作，
//在各个操作中解析请求再根据请求的操作找到对应的匿名函数,传入请求参数，将响
//应结果字符串返回给用户

var anonymousMap map[string]func(interface{}) string
var funcMap map[string]reflect.Type

func init() {
	funcMap = make(map[string]reflect.Type)

	anonymousMap = make(map[string]func(interface{}) string)

	structMap = make(map[string]reflect.Type)
}

func registerStruct(obj interface{}) {
	rv := reflect.TypeOf(obj)
	fmt.Println(rv.Name())
	structMap[rv.Name()] = rv
	fmt.Println(structMap[rv.Name()].NumField())
}

func registerFunc(funName string, fun interface{}) {
	funcMap[funName] = reflect.TypeOf(fun)
}

func registerAnonymousFunc(funName string, fun func(interface{}) string) {
	anonymousMap[funName] = fun
}
func registerRPCServer() error {
	return nil
}
