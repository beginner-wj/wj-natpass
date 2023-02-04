package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"wj-natpass/common"
)

var osName string

func main() {
	mainWindow()
}

func mainWindow() {
	a := app.New()
	w := a.NewWindow("wj-nat-pass")
	remoteAddress := widget.NewLabel("remote address:")
	remoteAddressTxt := widget.NewEntry()
	remoteAddressTxt.SetPlaceHolder("remote address,such as :1.1.1.1:8001")
	localhostAddress := widget.NewLabel("localhost address:")
	localhostAddressTxt := widget.NewEntry()
	localhostAddressTxt.SetPlaceHolder("localhost address,such as:127.0.0.1:8001") //默认127.0.0.1 也可以是局域网的端口
	//	v1 := container.NewBorder(layout.NewSpacer(), layout.NewSpacer(), remoteAddress, remoteAddressTxt)
	var submit *widget.Button
	cannelConn := false
	var connremote1, connremote2 net.Conn
	submit = widget.NewButton("connect", func() {
		if !cannelConn {
			currIp := checkParam(remoteAddressTxt, localhostAddressTxt)
			if !currIp {
				return
			}
			checkIpPortAvailable(localhostAddressTxt)
			submit.Disabled()
			common.StartForward()
			//连接
			go common.AddressToRemote(localhostAddressTxt.Text, remoteAddressTxt.Text, func(conn1 net.Conn, conn2 net.Conn) {
				cannelConn = true
				submit.Enable()
				connremote1 = conn1
				connremote2 = conn2
				submit.SetText("cancel connect")
				common.ToastSucc("connect succ")
			})

		} else {
			fmt.Println("=======取消连接=======")
			common.StopForward()
			connremote1.Close()
			connremote2.Close()
			cannelConn = false
			submit.SetText("connect")
		}
	})
	w.SetContent(container.NewVBox(remoteAddress, remoteAddressTxt, localhostAddress, localhostAddressTxt, submit))
	w.Resize(fyne.NewSize(300, 200))
	w.CenterOnScreen()
	w.ShowAndRun()
}

/*
***
检查端口是否可用
*/
func checkIpPortAvailable(localhostAddressTxt *widget.Entry) {
	port := localhostAddressTxt.Text
	cmdStr := "lsof -i:" + port
	bash := "/bin/bash"
	arg := "-c"
	if strings.ToUpper(osName) == "WINDOWS" {
		bash = "cmd.exe"
		cmdStr = "netstat -aon|findstr " + port
		arg = "/c"
	}
	//TODO 这里需要判断系统。如果是windows的话。还要另一种方式
	runCmd := exec.Command(bash, arg, cmdStr)
	if strings.ToUpper(osName) == "WINDOWS" {
		//本想直接代码里控制。但是mac没有这个属性。先去掉。build时隐藏即可
		//runCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := runCmd.Run()
	if err == nil {
		common.ToastSubmitFunc("tips", "port【"+port+"】is already use", func(res bool) {
		})
	}
}

func checkParam(remoteAddressTxt *widget.Entry, localAddressTxt *widget.Entry) bool {
	remoteTxt := remoteAddressTxt.Text
	ok, tip := checkIp(remoteTxt)
	if !ok {
		common.ToastError(tip)
		return ok
	}

	localTxt := localAddressTxt.Text
	ok, tip = checkIp(localTxt)
	if !ok {
		common.ToastError(tip)
		return ok
	}
	return true
}

func checkIp(address string) (bool, string) {
	ipAndPort := strings.Split(address, ":")
	if len(ipAndPort) != 2 {
		return false, "address【" + address + "】err，such as [ip:port]. "
	}
	ip := ipAndPort[0]
	port := ipAndPort[1]
	tip := checkPortNum(port)
	if tip != "" {
		return false, tip
	}
	pattern := `^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$`
	ok, err := regexp.MatchString(pattern, ip)
	if err != nil || !ok {
		return false, "ip error "
	}
	return ok, ""
}

func checkPortNum(port string) string {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return " port must by number : " + port
	}
	if portNum < 1 || portNum > 65535 {
		return " port 【" + port + "】must < 1 and > 65535 "
	}
	return ""
}
