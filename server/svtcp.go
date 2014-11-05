package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type PACKET struct {
	Ptype   uint8
	Content string
}

func send0(conn net.Conn, prog, location string) { //register
	msg := prog + "#" + MATHPROCEDURES + "#" + location

	enc := gob.NewEncoder(conn)
	packet := PACKET{0, msg}
	//fmt.Println("Sent:", packet)
	enc.Encode(&packet)
}

func send(conn net.Conn, packtype uint8, key string) {
	enc := gob.NewEncoder(conn)
	packet := PACKET{packtype, key}
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
