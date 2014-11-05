package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// generic packet structure
type PACKET struct {
	Ptype   uint8
	Content string
}

// send generic packet
func send(conn net.Conn, packtype uint8, key string) {
	enc := gob.NewEncoder(conn)
	packet := PACKET{packtype, key}
	fmt.Printf("")
	enc.Encode(&packet)
}

// get my ip address
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

// get port number
func getport(l net.Listener) string {
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}
