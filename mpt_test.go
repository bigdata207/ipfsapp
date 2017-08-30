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
}

//TestConverInfoNode 测试InfoNode与json字符串的互换
func TestConvertInfoNode(t *testing.T) {
	infoN := &InfoNode{}
	infoN.FileName = "test.txt"
	infoN.AbsolutePath, _ = NewPath("/home/vrit/test.txt")
	RP, _ := NewPath("vrit")
	m := make(map[string]string)
	m["0"] = "0"
	m["100"] = "100"
	var pies []piece
	pies = append(pies, NewPiece(1, "1", "1", []byte("1")))
	pies = append(pies, NewPiece(2, "2", "2", []byte("2")))
	infoN.Pieces = append(infoN.Pieces, pies...)
	infoN.UploadRoot = parent{Hash: "qbvfdhsvfgh", RelativePath: RP}
	data, err := InfoNode2Json(infoN)
	if err == nil {
		fmt.Println(string(data))
		infoN2, err := NewInfoNodeFromJson(data)
		if err == nil {
			fmt.Println(infoN2.AbsolutePath)

		}
	}
}
