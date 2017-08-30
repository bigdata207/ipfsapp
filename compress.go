package ipfsapp

import (
	_ "bufio"
	"errors"
	_ "io/ioutil"
	_ "os"
	_ "path/filepath"
)

func zipCompress(ifname, ofname string) error {
	return nil
}
func tarCompress(ifname, ofname string) error {
	return nil
}
func targzCompress(ifname, ofname string) error {
	return nil
}
func tarxzCompress(ifname, ofname string) error {
	return nil
}
func sevenzCompress(ifname, ofname string) error {
	return nil
}
func rarCompress(ifname, ofname string) error {
	return nil
}

//Compress 文件压缩接口，需要输入输出文件名
func Compress(ifname, mode string, ofname ...string) error {
	var err error
	if len(ofname) == 0 {
		ofname = append(ofname, ifname+"."+mode)
	}
	switch mode {
	case "zip":
		err = zipCompress(ifname, ofname[0])
	case "tar":
		err = tarCompress(ifname, ofname[0])
	case "tar.gz":
		err = targzCompress(ifname, ofname[0])
	case "tar.xz":
		err = tarxzCompress(ifname, ofname[0])
	case "7z":
		err = sevenzCompress(ifname, ofname[0])
	case "rar":
		err = rarCompress(ifname, ofname[0])
	default:
		err = errors.New("Don't support this compress formate")
	}
	return err
}
