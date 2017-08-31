package ipfsapp

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func testIPFS() bool {
	cmd := exec.Command("whereis", "ipfs")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
	return len(strings.Split(out.String(), ":")) > 1
}

func installIPFS(installPath ...string) error {
	downoadLink := "https://dist.ipfs.io/go-ipfs/v0.4.10/go-ipfs_v0.4.10_linux-amd64.tar.gz"
	tmpDir := "/tmp/ipfs"
	path, _ := os.Getwd()
	if len(installPath) > 0 {
		path = installPath[0]
	}
	//parentFolder, file := filepath.Split(path)
	f, err := os.Open(path)
	if err != nil {
		os.MkdirAll(path, os.ModeDir)
		f, _ = os.Open(path)
	} else {
		//-1 表示读取所有目录
		dirs, err := f.Readdir(-1)
		if err == nil {
			if len(dirs) > 0 {
				if strend(path) == '/' {
					path += "ipfs"
				} else {
					path += "/ipfs"
				}
				err = os.Mkdir(path, os.ModeDir)
			}
		}
		f, _ = os.Open(path)
	}
	fname := strings.Split(downoadLink, "/")[len(strings.Split(downoadLink, "/"))-1]
	cmd := exec.Command(fmt.Sprintf("mkdir -P %v && wget -P %v %v && tar -zcvf %v -C %v", tmpDir, tmpDir, downoadLink, tmpDir+"/"+fname, path))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
	u, _ := user.Current()
	bin := path + "/go-ipfs"
	cmd = exec.Command(fmt.Sprintf("echo \"export PATH=$PATH:%v\" | tee -a /home/%v/.bashrc echo \"export PATH=$PATH:%v\" | tee -a /home/%v/.zshrc", bin, u.HomeDir, bin, u.HomeDir))
	err = cmd.Run()
	return err
}

func ipfsInit() {
	cmd := exec.Command("ipfs", "init")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsStart() {
	cmd := exec.Command("ipfs", "daemon")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsGet(ipfsHash string) {
	cmd := exec.Command("ipfs", "get", ipfsHash)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsAddFolder(dir string) {
	cmd := exec.Command("ipfs add", "-r", dir)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsAdd(fname string) {
	cmd := exec.Command("ipfs", "add", fname)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func ipfsStats(arg string) {
	args := []string{"bitswap", "bw", "repo"}
	b := ContainEle(strArr2InterArr(args), arg)
	if b == -1 {
		return
	}
	cmd := exec.Command("ipfs", "stat", arg)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}
