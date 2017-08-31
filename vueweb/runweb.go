package ipfsapp

import (
	"errors"
	"os/exec"
	curUser "os/user"
	"time"
)

var webRoot string

func init() {
	webRoot = "vueweb"
}

func switchNpmRegistery() error {
	user, err := curUser.Current()
	if err != nil {
		return errors.New("Can't get user home dir")
	}
	cmd := exec.Command("echo \"registry=https://registry.npm.taobao.org\" ", ">", user.HomeDir)
	err = cmd.Run()
	return err
}

func npmInstall(c chan error) {
	cmd := exec.Command("cd", webRoot, "&& npm install")
	err := cmd.Run()
	c <- err
}

func npmStart(c chan error) {
	cmd := exec.Command("cd", webRoot, "&& npm satrt")
	err := cmd.Start()
	if err != nil {
		c <- err
	}
	time.Sleep(time.Second * 3)
	if cmd.ProcessState.Exited() {
		c <- errors.New("Web 启动失败")
	}
}

func npmRunDev(c chan error) {
	cmd := exec.Command("cd", webRoot, "&& npm run dev")
	err := cmd.Start()
	if err != nil {
		c <- err
	}
	time.Sleep(time.Second * 3)
	if cmd.ProcessState.Exited() {
		c <- errors.New("Web 启动失败")
	}
}
