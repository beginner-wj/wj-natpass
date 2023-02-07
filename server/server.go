package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"wj-natpass/common"
)

func main() {
	printWelcome()
	args := os.Args
	cmdlength := len(args)
	if cmdlength <= 2 {
		printHelp()
		os.Exit(0)
	}
	tcpBusiness(args, cmdlength)
}

func tcpBusiness(args []string, cmdlength int) {
	switch args[1] {
	case "-listen":
		if cmdlength < 3 {
			log.Fatalln(`-listen need two arguments, like "main -listen 1997 2017".`)
		}
		port1 := checkPort(args[2])
		port2 := checkPort(args[3])
		portToPort(port1, port2)
		break
	case "-tran":
		if cmdlength < 3 {
			log.Fatalln(`-listen need two arguments, like "main -tran 8081 192.168.31.2:8081".`)
		}
		port := checkPort(args[2])
		var remoteAddress string
		if common.CheckIp(args[3]) {
			remoteAddress = args[3]
		}
		split := strings.SplitN(remoteAddress, ":", 2)
		log.Println("[√]", "start to transmit address:", remoteAddress, "to address:", split[0]+":"+port)
		portToRemote(port, remoteAddress)
		break
	case "-slave":
		if cmdlength < 3 {
			log.Fatalln(`-slave need two arguments, like "wj-nat -slave 127.0.0.1:3389 8.8.8.8:1997".`)
		}
		var remoteAddress string
		if common.CheckIp(args[2]) {
			remoteAddress = args[2]
		}
		var remoteAddress1 string
		if common.CheckIp(args[3]) {
			remoteAddress1 = args[3]
		}
		log.Println("[√]", "start to connect address:", remoteAddress, "and address:", remoteAddress1)
		AddressToRemote(remoteAddress, remoteAddress1)
		break
	}
}

const timeout = 5

func AddressToRemote(remoteAddress string, remoteAddress1 string) {
	for {
		log.Println("[+]", "try to connect host:["+remoteAddress+"] and ["+remoteAddress1+"]")
		var connremote1, connremote2 net.Conn
		var err error
		for {
			connremote1, err = net.Dial("tcp", remoteAddress)
			if err == nil {
				log.Println("[→]", "connect ["+remoteAddress+"] success.")
				break
			}
		}
		for {
			connremote2, err = net.Dial("tcp", remoteAddress1)
			if err == nil {
				log.Println("[→]", "connect ["+remoteAddress1+"] success.")
				break
			}
		}
		common.Forward(connremote1, connremote2)
	}
}

func portToRemote(port string, remoteAddress string) {
	listener1 := listen_port("0.0.0.0:" + port)
	for {
		conn := common.Accept(listener1)
		if conn != nil {
			log.Println("[+]", "start connect host:["+remoteAddress+"]")
			conn.Close()
			log.Println("[←]", "close the connect at local:["+conn.LocalAddr().String()+"] and remote:["+conn.RemoteAddr().String()+"]")
			time.Sleep(timeout * time.Second)
			return
		}
		go connectRemoteAddressAndforward(conn, remoteAddress)
	}
}

func connectRemoteAddressAndforward(conn net.Conn, remoteAddress string) {
	connremote, err := net.Dial("tcp", remoteAddress)
	if err != nil {
		log.Println("[x ]", "start connect host:["+remoteAddress+"] fail :"+err.Error())
		panic("start connect host:[" + remoteAddress + "] fail :" + err.Error())
	}
	common.Forward(connremote, conn)
}

func portToPort(port1 string, port2 string) {
	listener1 := listen_port("0.0.0.0:" + port1)
	listener2 := listen_port("0.0.0.0:" + port2)
	for {
		conn1 := common.Accept(listener1)
		conn2 := common.Accept(listener2)
		if conn1 == nil || conn2 == nil {
			log.Println("[x]", "common.Accept client faild. retry in ", timeout, " seconds. ")
			time.Sleep(timeout * time.Second)
			continue
		}
		common.Forward(conn1, conn2)

		//defer conn1.Close()
		//defer conn2.Close()
	}
}

func listen_port(port string) net.Listener {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(" listen  : ", port, " fail: ", err.Error())
		panic(" listen  : " + port + " fail: " + err.Error())
	}
	return listener
}

func checkPort(port string) string {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalln(" port must by number : ", port)
		panic(" port must by number : " + port)
	}
	if portNum < 1 || portNum > 65535 {
		log.Fatalln(" port must < 1 and > 65535 ")
		panic(" port must < 1 and > 65535 ")
	}
	return port
}

func printWelcome() {
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println("| Welcome to use wj-nat |")
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println()
	time.Sleep(time.Second)
}
func printHelp() {
	fmt.Println(`usage: "-listen port1 port2" example: "wj-nat -listen 1997 2017" `)
	fmt.Println(`       "-tran port1 ip:port2" example: "wj-nat -tran 1997 192.168.1.2:3389" `)
	fmt.Println(`       "-slave ip1:port1 ip2:port2" example: "wj-nat -slave 127.0.0.1:3389 8.8.8.8:1997" `)
	fmt.Println(`============================================================`)
	fmt.Println(`optional argument: "-log logpath" . example: "wj-nat -listen 1997 2017 -log d:/nb" `)
	fmt.Println(`log filename format: Y_m_d_H_i_s-agrs1-args2-args3.log`)
	fmt.Println(`============================================================`)
	fmt.Println(`if you want more help, please read "README.md". `)
}
