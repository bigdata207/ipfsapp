package ipfsapp

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	exts       = []string{"zip", "rar", "7z", "tar", "tar.gz", "tar.xz", "tar.bz2"}
	skipfolder = false
)

func zipCompress(ifpath, ofpath string, skipfolder ...bool) error {
	return nil
}
func tarCompress(ifpath, ofpath string, skipfolder ...bool) error {
	return nil
}
func tarbz2Compress(ifpath, ofpath string, skipfolder ...bool) error {
	return nil
}

func targzCompress(ifpath, ofpath string, skipfolder ...bool) error {
	// file write
	folder, _ := filepath.Split(ofpath)
	os.MkdirAll(folder, os.ModeDir)
	fw, err := os.Create(ofpath)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 打开文件夹
	dir, err := os.Open(ifpath)
	if err != nil {
		panic(nil)
	}
	defer dir.Close()

	// 读取文件列表
	//fis, err := dir.Readdir(0)
	fis, err := ioutil.ReadDir(ifpath)

	if len(fis) == 0 {
		p, fn := filepath.Split(ifpath)
		os.Mkdir(p+"tmp", os.ModeDir)
		cmd := exec.Command("mv", ifpath, p+"tmp")
		err = cmd.Run()
		ifpath = p + "tmp/" + fn
		fis, err = ioutil.ReadDir(ifpath)
		defer func() {
			cmd := exec.Command("mv", ifpath, p)
			err = cmd.Run()
			os.RemoveAll(p + "tmp")
		}()
	} else {
		// 遍历文件列表
		for _, fi := range fis {
			// 逃过文件夹, 我这里就不递归了
			if len(skipfolder) > 0 && skipfolder[0] && fi.IsDir() {
				continue
			}

			// 打印文件名称
			fmt.Println(fi.Name())

			// 打开文件
			fr, err := os.Open(dir.Name() + "/" + fi.Name())
			if err != nil {
				panic(err)
			}
			defer fr.Close()

			// 信息头
			h := new(tar.Header)
			h.Name = fi.Name()
			h.Size = fi.Size()
			h.Mode = int64(fi.Mode())
			h.ModTime = fi.ModTime()

			// 写信息头
			err = tw.WriteHeader(h)
			if err != nil {
				panic(err)
			}

			// 写文件
			_, err = io.Copy(tw, fr)
			if err != nil {
				panic(err)
			}
		}
	}
	fmt.Println("tar.gz ok")
	return err
}
func tarxzCompress(ifpath, ofpath string, skipfolder ...bool) error {
	return nil
}
func sevenzCompress(ifpath, ofpath string, skipfolder ...bool) error {
	return nil
}
func rarCompress(ifpath, ofpath string, skipfolder ...bool) error {
	return nil
}

//Compress 文件压缩接口，需要输入输出文件名
func Compress(ifpath, formate string, ofpath ...string) error {
	var err error
	if len(ofpath) == 0 {
		ofpath = append(ofpath, ifpath+"."+formate)
	}
	switch formate {
	case "zip":
		err = zipCompress(ifpath, ofpath[0])
	case "tar":
		err = tarCompress(ifpath, ofpath[0])
	case "tar.gz":
		err = targzCompress(ifpath, ofpath[0])
	case "tar.xz":
		err = tarxzCompress(ifpath, ofpath[0])
	case "7z":
		err = sevenzCompress(ifpath, ofpath[0])
	case "rar":
		err = rarCompress(ifpath, ofpath[0])
	case "tar.bz2":
		tarbz2Compress(ifpath, ofpath[0])
	default:
		err = errors.New("Don't support this compress formate")
	}
	return err
}
