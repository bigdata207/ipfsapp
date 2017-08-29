package main

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
	/**
	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi)
	fo, err := os.Create("2.cry")
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(fo)
	bufsize := 32
	buf := make([]byte, bufsize)
	padsize := 0
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		if n == bufsize {
			data, err = AesEncrypt(buf[:n], key, false)
		} else {
			data, err = AesEncrypt(buf[:n], key, true)
			pad = n
		}
		if err != nil {
			panic(err)
		}
		if n2, err := w.Write(data); err != nil {
			panic(err)
		} else if n2 != len(data) {
			panic("error in writing")
		}
	}
	padtext := bytes.Repeat([]byte{0}, padsize)
	if err = w.Flush(); err != nil {
		panic(err)
	}
	fo.Close()

	rfi, err := os.Open("2.cry")
	if err != nil {
		panic(err)
	}
	defer rfi.Close()
	r = bufio.NewReader(rfi)
	rfo, err := os.Create("2.mp3")
	if err != nil {
		panic(err)
	}
	defer rfo.Close()
	rw := bufio.NewWriter(rfo)

	buf = make([]byte, 128)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		origData, err := AesDecrypt(buf[:n], key, falsel)
		if err != nil {
			panic(err)
		}
		if n2, err := rw.Write(origData); err != nil {
			panic(err)
		} else if n2 != len(origData) {
			panic("error in writing")
		}
	}

	if err = rw.Flush(); err != nil {
		panic(err)
	}
	**/
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
