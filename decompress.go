package ipfsapp

import (
	_ "bufio"
	"errors"
	_ "io/ioutil"
	_ "os"
	_ "path/filepath"
)

func zipUnCompress(ifname, ofname string) error {
	return nil
}
func tarUnCompress(ifname, ofname string) error {
	return nil
}
func targzUnCompress(ifname, ofname string) error {
	return nil
}
func tarxzUnCompress(ifname, ofname string) error {
	return nil
}
func sevenzUnCompress(ifname, ofname string) error {
	return nil
}
func rarUnCompress(ifname, ofname string) error {
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
