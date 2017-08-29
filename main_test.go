package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
	fmt.Println("I am main")
	curPath, _ := os.Getwd()
	fmt.Println(curPath)
	fmt.Println(filepath.Split(curPath))
	arr := []int{1, 2, 3}
	arr = append(arr, 1)
	fmt.Println(containEle(intArr2InterArr(arr), 2))
}
