package ipfsapp

import "github.com/andlabs/ui"

//LoginBox 登录界面
func LoginBox() {
	user := ui.NewEntry()
	pass := ui.NewEntry()
	button := ui.NewButton("Login")
	res := ui.NewLabel("")
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Login UI"), false)
	box.Append(user, false)
	box.Append(pass, false)
	box.Append(button, false)
	box.Append(res, false)
	window := ui.NewWindow("Hello", 200, 200, false)
	window.SetChild(box)
	button.OnClicked(func(*ui.Button) {
		res.SetText("Hello, " + user.Text() + "@" + pass.Text() + "!")
	})
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}

//RegisterBox 注册界面
func RegisterBox() {}

//AlertBox 警告弹窗
func AlertBox(warnMess string) {}

//ConfirmBox 确认弹窗
func ConfirmBox() {}

//ListBox 列表展示界面
func ListBox() {}

//ManagerBox 管理界面
func ManagerBox() {}

//MainUI 主界面
func MainUI() {}
