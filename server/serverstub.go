// serverstub
package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const MATHPROCEDURES = "Min$Max$Multiply$Sort"

var myip string
var myprogramver string

func handshake(TCPconn net.Conn) {
	//negotiate with the client
	dec := gob.NewDecoder(TCPconn)
	packet := PACKET{} //Math1#Multiply#2#2x3#3x1#2x1
	dec.Decode(&packet)
	if packet.Ptype != 4 {
		return
	}

	//Math1#Multiply#2
	req := strings.Split(packet.Content, "#")
	fmt.Println(req)

	// check if is provided
	if req[0] != myprogramver || provided(req[1], req[2]) == false {
		fmt.Println("Service not provided")
		return
	}
	fmt.Println("====================================================================\n[[[CLIENT REQUESTED SERVICE]]]", req[0], req[1])
	addr := net.UDPAddr{
		Port: 0,
		IP:   net.ParseIP("0.0.0.0"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	fmt.Println("====================================================================\n[[[UDP CONNECTION ESTABLISHED FOR " + req[1] + "]]]")

	printerr(err)

	tmp := ser.LocalAddr().String()
	add := strings.Split(tmp, ":")

	// tell the client UDP address
	location := myip + ":" + add[1]

	//send client the UDP address
	send(TCPconn, 5, location)
	for {
		handleUDP(ser, req)
	}
}

func conn2pm() net.Conn {
	//client stub reading port mapper location
	loc, _ := ioutil.ReadFile("../pmlocation")

	//connect
	conn, _ := net.Dial("tcp", string(loc))
	if conn == nil {
		p("Cannot connect to port mapper! Very likely port mapper isn't started.")
		os.Exit(1)
	}
	return conn
}

func startDaemon(pm net.Conn, progver string) {
	//create new thread listening to client
	l, _ := net.Listen("tcp", ":0")

	myip = getmyip()
	port := getport(l)
	location := myip + ":" + port

	//send register packet
	send0(pm, progver, location)

	//wait for ack
	dec := gob.NewDecoder(pm)
	packet := PACKET{}
	dec.Decode(&packet)
	//fmt.Println(packet)

	if packet.Content != progver {
		p("Register unsuccessful!")
		return
	}
	pm.Close() // finished with port mapper

	for { // busy waiting for connection
		clientconn, _ := l.Accept() // this blocks until connection or error
		go handshake(clientconn)    // a goroutine handles conn so that the loop can accept other connections
	}
}

func Register(progver string) {
	myprogramver = progver
	pm := conn2pm()
	startDaemon(pm, progver)
	pm.Close()
	return
}

func provided(name, dimension string) bool {
	ret := false
	switch name {
	case "Min", "Max", "Sort":
		if dimension == "1" {
			ret = true
		}
	case "Multiply":
		if dimension == "2" {
			ret = true
		}
	default:
	}
	return ret
}
