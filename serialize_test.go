package ipfsapp

import (
	"fmt"
	"testing"
)

func TestSerizlize(t *testing.T) {
	fmt.Println("Test Serialize")
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
	data, err := struct2json(infoN)
	if err == nil {
		fmt.Println(string(data))
		infoN2 := &InfoNode{}
		json2struct(data, infoN2)
		fmt.Println(infoN2.Pieces)
	}
}
