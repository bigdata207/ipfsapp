package ipfsapp

import (
	"archive/tar"
	_ "bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	_ "io/ioutil"
	"os"
	"path/filepath"
	_ "path/filepath"
)

func zipUnCompress(ifpath, ofpath string) error {
	return nil
}
func tarUnCompress(ifpath, ofpath string) error {
	return nil
}
func targzUnCompress(ifpath, ofpath string) error {
	_, fn := filepath.Split(ifpath)

	if strend(ofpath) == '/' {
		ofpath += fn[:(len(fn) - 7)]
	} else {
		ofpath += ("/" + fn[:(len(fn)-7)])
	}

	fr, err := os.Open(ifpath)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// tar read
	tr := tar.NewReader(gr)

	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// 显示文件
		fmt.Println(h.Name)

		// 打开文件
		fw, err := os.OpenFile(ofpath+h.Name, os.O_CREATE|os.O_WRONLY, 0644 /*os.FileMode(h.Mode)*/)
		if err != nil {
			panic(err)
		}
		defer fw.Close()

		// 写文件
		_, err = io.Copy(fw, tr)
		if err != nil {
			panic(err)
		}

	}

	fmt.Println("un tar.gz ok")
	return nil
}
func tarxzUnCompress(ifpath, ofpath string) error {
	return nil
}
func sevenzUnCompress(ifpath, ofpath string) error {
	return nil
}
func rarUnCompress(ifpath, ofpath string) error {
	return nil
}

//UnCompress 文件解压缩接口
func UnCompress(ifpath, mode string, ofpath ...string) (err error) {

	if len(ofpath) == 0 {
		//dir, ifp := filepath.Split(ifpath)
		ofpath = append(ofpath, ifpath+"."+mode)
	} else {
		if ofpath[0][(len(ofpath[0])-len(mode)):] != mode {
			err = errors.New("format doesn't match mode")
			return
		}
	}
	switch mode {
	case "zip":
		err = zipUnCompress(ifpath, ofpath[0])
	case "tar":
		err = tarUnCompress(ifpath, ofpath[0])
	case "tar.gz":
		err = targzUnCompress(ifpath, ofpath[0])
	case "tar.xz":
		err = tarxzCompress(ifpath, ofpath[0])
	case "7z":
		err = sevenzUnCompress(ifpath, ofpath[0])
	case "rar":
		err = rarUnCompress(ifpath, ofpath[0])
	default:
		err = errors.New("Don't support this compress formate")
	}
	return
}
