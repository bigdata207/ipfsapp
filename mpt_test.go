package ipfsapp

import (
	"fmt"
	"testing"
)

func TestMPT(t *testing.T) {
	fmt.Println("Test MPT")
	calSalt()
	p, err := NewPath("/home")
	if err == nil {
		fmt.Println(p.Add("vrit"))
		fmt.Println(p.Back("root"))
		fmt.Println(p.Sub("/"))
	} else {
		fmt.Println(err)
	}
	in := IndexNode{nType: FileNode{}}
	in.nType = DirNode{}
	fmt.Println(in.Type())

	fmt.Println(p.KeyHash())
	fmt.Println(len(p.KeyHash()))
	p, _ = p.Add("go")
	fmt.Println(p.KeyHash())
	fmt.Println(Get("h"))
}
