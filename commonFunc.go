package ipfsapp

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"strings"
)

//ReadLine 从终端读取一行输入
func ReadLine() ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	data, b, err := reader.ReadLine()
	fmt.Println(b)
	return data, err
}

//strend　获取字符串的最后一个字符
func strend(s string) byte {
	return s[len(s)-1]
}

//func2anonymousFunc　简单地将一个已有的函数转化成一个匿名函数
func func2anonymousFunc(funcName string) func(...interface{}) interface{} {
	//f := reflect.TypeOf(strend)
	f := funcMap[funcName]
	fmt.Println("func para num: ", f.NumIn())
	for i := 0; i < f.NumIn(); i++ {
		fmt.Printf("in %d: %v\n", i, f.In(i).Name())
	}
	fmt.Println("func return num:", f.NumOut())
	for i := 0; i < f.NumOut(); i++ {
		fmt.Printf("out %d: %v\n", i, f.Out(i).Name())
	}
	fun := func(a ...interface{}) interface{} {
		if len(a) == f.NumIn() {
			return interface{}(strend(a[0].(string)))
		}
		return nil
	}
	return fun
}

//NewCounter 返回一个从begin开始的计数器函数
func NewCounter(begin int) func() int {
	i := begin
	return func() int {
		fmt.Println(i)
		t := i
		i++
		return t
	}
}

//RelaPath2AbsoPath 相对路径转绝对路径
func RelaPath2AbsoPath(relaPath string) string {
	curDir, _ := os.Getwd()
	if relaPath[0] == '.' {
		return curDir + "/" + relaPath[1:]
	} else if relaPath[0] != '/' {
		return curDir + "/" + relaPath
	}
	return relaPath
}

//Search 从给定的IndexNode开始查找和key匹配的路径
func Search(in IndexNode, key string) []Path {
	paths := make([]Path, 0)
	if strings.Contains(in.kv.k, key) {
		paths = append(paths, in.kv.v)
	}

	if in.IsDir() {
		for _, v := range in.childBlocks {
			tmp := Search(v, key)
			if len(tmp) == 0 {
				return paths
			}
			paths = append(paths, tmp...)
		}
	}
	return paths
}

func dfs(kv kvNode) map[string]IndexNode {
	dirPath := kv.v
	ins := make(map[string]IndexNode)
	childs, err := ioutil.ReadDir(dirPath.String())
	if err == nil {
		for _, v := range childs {
			tmp := deepth(dirPath.Add(v.Name()))
			diff := tmp.Sub(dirPath)
			if len(diff) > 0 {
				ins[diff] = NewIndexNode(DirNode{}, kvNode{k: diff, v: tmp})
			}
		}
	}
	return ins
}

func deepth(nodePath Path, e ...interface{}) Path {
	childs, err := ioutil.ReadDir(nodePath.String())
	if err == nil {
		if len(childs) == 1 && childs[0].IsDir() {
			childPath, err := nodePath.Add(childs[0].Name())
			if err == nil {
				return deepth(childPath)
			}
		}
	}
	return ""
}
func intArr2InterArr(arr []int) []interface{} {
	tmp := make([]interface{}, len(arr), len(arr))
	for i, v := range arr {
		tmp[i] = interface{}(v)
	}
	return tmp
}
func strArr2InterArr(arr []string) []interface{} {
	tmp := make([]interface{}, len(arr), len(arr))
	for i, v := range arr {
		tmp[i] = interface{}(v)
	}
	return tmp
}

//ContainEle 判断数组中是否含有某个元素
func ContainEle(arr []interface{}, ele interface{}) int {
	if len(arr) > 0 {
		if strings.Compare(reflect.TypeOf(arr[0]).Name(), reflect.TypeOf(ele).Name()) == 0 {
			for i, v := range arr {
				if v == ele {
					return i
				}
			}
		}
	}
	return -1
}

//LinesCounter 统计文件下go源码行数
func LinesCounter(dir string, suffix ...string) int {
	if len(suffix) == 0 {
		suffix = append(suffix, "go")
	}
	lines := 0
	fs, err := ioutil.ReadDir(dir)
	if err == nil {
		for _, v := range fs {
			if v.IsDir() {
				lines += LinesCounter(dir + "/" + v.Name())
			} else {
				s := strings.Split(v.Name(), ".")
				if ContainEle(strArr2InterArr(suffix), s[len(s)-1]) > -1 {
					f, _ := ioutil.ReadFile(dir + "/" + v.Name())
					lines += len(strings.Split(string(f), "\n"))
				}
			}
		}
	}
	return lines
}

func getIP() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ips := make([]string, 0)
	for _, address := range addrs {

		// ipnet.IP.IsLoopback()检查ip地址判断是否回环地址
		//if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
		if ipnet, ok := address.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String())
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	return ips
}
