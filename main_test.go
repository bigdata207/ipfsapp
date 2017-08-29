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
}
