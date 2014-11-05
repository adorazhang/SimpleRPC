package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type pac struct {
	Ptype   uint8
	Content string
}

func send3(conn net.Conn, data string) { //return server location
	enc := gob.NewEncoder(conn)
	packet := pac{3, data}
	enc.Encode(&packet)
}

func send1(conn net.Conn, data string) { //send ack to server
	enc := gob.NewEncoder(conn)
	packet := pac{1, data}
	enc.Encode(&packet)
}

func getmyip() string {
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		tip := address.(*net.IPAddr).IP.String()
		if tip != "0.0.0.0" && !strings.HasPrefix(tip, "192.168.") && !strings.HasPrefix(tip, "10.") {
			return tip
		}
	}
	fmt.Println("Can't get local IP!")
	os.Exit(1)
	return ""
}

func getport(l net.Listener) string {
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}
