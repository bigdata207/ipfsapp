package ipfsapp

import (
	"fmt"
	"os"
	"os/user"
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
	u, _ := user.Current()
	fmt.Println(u.HomeDir)
}
