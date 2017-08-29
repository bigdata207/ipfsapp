package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
)

var decrypted string

//var privateKey []byte
//var publicKey []byte

func rK() {
	flag.StringVar(&decrypted, "d", "", "加密过的数据")
	flag.Parse()
	privateKey, _ = ioutil.ReadFile("private.pem")

	publicKey, _ = ioutil.ReadFile("public.pem")
}

func tryRSA() {
	//InitKeys(&publicKey, &privateKey)
	var data []byte
	var err error
	if decrypted != "" {
		data, err = base64.StdEncoding.DecodeString(decrypted)
		if err != nil {
			panic(err)
		}
	} else {
		args := flag.Args()
		fileName := args[0]
		//f, err := ioutil.ReadFile(fileName)
		//fmt.Println(len(f))
		//ff := f[0:128]
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
		defer fo.Close()
		w := bufio.NewWriter(fo)

		buf := make([]byte, 128)
		for {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}
			data, err = RsaEncrypt(buf[:n])
			if err != nil {
				panic(err)
			}
			if n2, err := w.Write(data); err != nil {
				panic(err)
			} else if n2 != len(data) {
				panic("error in writing")
			}
		}

		if err = w.Flush(); err != nil {
			panic(err)
		}

		//fmt.Println("rsa encrypt base64:" + base64.StdEncoding.EncodeToString(data))
	}
	/**
		origData, err := RsaDecrypt(data)
		fmt.Printf("%d : %d\n", len(data), len(origData))
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(origData))
	**/

	rfi, err := os.Open("2.cry")
	if err != nil {
		panic(err)
	}
	defer rfi.Close()
	r := bufio.NewReader(rfi)
	rfo, err := os.Create("2.mp3")
	if err != nil {
		panic(err)
	}
	defer rfo.Close()
	rw := bufio.NewWriter(rfo)

	buf := make([]byte, 256)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		origData, err := RsaDecrypt(buf[:n])
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

}

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	//fmt.Println(string(publicKey))
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	//fmt.Println(string(privateKey))
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
