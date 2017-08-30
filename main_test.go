package ipfsapp

import (
	"fmt"
	"github.com/mackzhong/ipfsapp/restful"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
	fmt.Println("I am ipfsapp")
	curPath, _ := os.Getwd()
	fmt.Println(curPath)
	fmt.Println(filepath.Split(curPath))
	arr := []int{1, 2, 3}
	arr = append(arr, 1)
	fmt.Println(ContainEle(intArr2InterArr(arr), 2))
	fmt.Println(testIPFS())
	a := "1234"
	fmt.Println(string(strend(a)))
	fmt.Println(restful.Add(1, 2))
}
