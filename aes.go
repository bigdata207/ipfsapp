package ipfsapp

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	_ "crypto/rand"
	_ "encoding/hex"
	"flag"
	"fmt"
	_ "io"
	"io/ioutil"
	"os"
	"time"
)

var data []byte
var err error
var fileName string
var bufsize int
var key []byte

func readKey() {
	flag.Parse()

	publicKey, _ := ioutil.ReadFile("public.pem")
	key = publicKey[:32]
	//fmt.Println(key)
	args := flag.Args()
	fileName = args[0]
	bufsize = 128
}

func tryAES() {
	t1 := time.Now()
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	padsize := bufsize - len(data)%bufsize
	//fmt.Println("data padsize: ", padsize)
	padtext := bytes.Repeat([]byte{byte(padsize)}, padsize)
	data = append(data, padtext...)
	//fmt.Println(int(data[len(data)-1]))
	fo, err := os.Create("2.cry")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(fo)
	for i := 0; i < len(data); i += bufsize {
		crypted, err := AesEncrypt(data[i:i+bufsize], key)
		if err != nil {
			panic(err)
		}
		if n2, err := w.Write(crypted); err != nil {
			panic(err)
		} else if n2 != bufsize {
			panic("error in writing")
		}
	}
	if err = w.Flush(); err != nil {
		panic(err)
	}
	fo.Close()

	cryptodata, err := ioutil.ReadFile("2.cry")
	if err != nil {
		panic(err)
	}
	//fmt.Println(len(cryptodata) % bufsize)
	fo, err = os.Create("2.mp3")
	if err != nil {
		panic(err)
	}
	w = bufio.NewWriter(fo)
	i := 0
	for ; i < len(cryptodata)-bufsize; i += bufsize {
		orig, err := AesDecrypt(cryptodata[i:i+bufsize], key)
		if err != nil {
			panic(err)
		}
		if n2, err := w.Write(orig); err != nil {
			panic(err)
		} else if n2 != bufsize {
			panic("error in writing")
		}
	}
	orig, err := AesDecrypt(cryptodata[i:len(cryptodata)], key)
	w.Write(orig[:(len(orig) - int(orig[len(orig)-1]))])
	if err = w.Flush(); err != nil {
		panic(err)
	}
	fo.Close()
	rt := time.Since(t1)
	fmt.Println("Runnning Time: ", rt)
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	//fmt.Println("De:", blockSize)
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	//fmt.Println(len(origData))

	//origData = PKCS5UnPadding(origData)

	//fmt.Println("orig:", len(origData))
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
