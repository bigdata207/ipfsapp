package ipfsapp

import (
	"fmt"
	"testing"
)

func TestTransFile(t *testing.T) {
	fmt.Println("测试多线程传输文件")
	fmt.Println(goTransFile("1.mp3"))
}
