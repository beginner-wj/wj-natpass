package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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
		if CheckIp(args[3]) {
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
		if CheckIp(args[2]) {
			remoteAddress = args[2]
		}
		var remoteAddress1 string
		if CheckIp(args[3]) {
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
		forward(connremote1, connremote2)
	}
}

func portToRemote(port string, remoteAddress string) {
	listener1 := listen_port("0.0.0.0:" + port)
	for {
		conn := accept(listener1)
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
	forward(connremote, conn)
}

func portToPort(port1 string, port2 string) {
	listener1 := listen_port("0.0.0.0:" + port1)
	listener2 := listen_port("0.0.0.0:" + port2)
	for {
		conn1 := accept(listener1)
		conn2 := accept(listener2)
		if conn1 == nil || conn2 == nil {
			log.Println("[x]", "accept client faild. retry in ", timeout, " seconds. ")
			time.Sleep(timeout * time.Second)
			continue
		}
		forward(conn1, conn2)

		//defer conn1.Close()
		//defer conn2.Close()
	}
}

func forward(conn1 net.Conn, conn2 net.Conn) {
	log.Printf("[+] start transmit. [%s],[%s] <-> [%s],[%s] \n", conn1.LocalAddr().String(), conn1.RemoteAddr().String(), conn2.LocalAddr().String(), conn2.RemoteAddr().String())
	var wg sync.WaitGroup
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	wg.Wait()
}

func connCopy(conn1 net.Conn, conn2 net.Conn, group *sync.WaitGroup) {
	io.Copy(conn1, conn2)
	conn1.Close()
	group.Done()
}

func accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		panic("listen " + conn.LocalAddr().String() + " fail " + err.Error())
	}
	return conn
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
	fmt.Println(`usage: "-listen port1 port2" example: "nb -listen 1997 2017" `)
	fmt.Println(`       "-tran port1 ip:port2" example: "nb -tran 1997 192.168.1.2:3389" `)
	fmt.Println(`       "-slave ip1:port1 ip2:port2" example: "nb -slave 127.0.0.1:3389 8.8.8.8:1997" `)
	fmt.Println(`============================================================`)
	fmt.Println(`optional argument: "-log logpath" . example: "nb -listen 1997 2017 -log d:/nb" `)
	fmt.Println(`log filename format: Y_m_d_H_i_s-agrs1-args2-args3.log`)
	fmt.Println(`============================================================`)
	fmt.Println(`if you want more help, please read "README.md". `)
}

func CheckIp(address string) bool {
	ipAndPort := strings.Split(address, ":")
	if len(ipAndPort) != 2 {
		log.Fatalln("[x]", "address error. should be a string like [ip:port]. ")
	}
	ip := ipAndPort[0]
	port := ipAndPort[1]
	checkPort(port)
	pattern := `^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$`
	ok, err := regexp.MatchString(pattern, ip)
	if err != nil || !ok {
		log.Fatalln("[x]", "ip error. ")
	}
	return ok
}
