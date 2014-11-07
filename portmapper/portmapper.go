// portmapper
// portmapper will always use TCP connections with the clients and servers
package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

/* table looks like
"Math1Multiply" -> "123.12.1.123:1111"
"Math2Min" -> "222.222.222.222:1212" */
var table map[string][]string

/* store current pointer */
var cur map[string]int

func addEntry(key, value string) {
	table[key] = append(table[key], value)
	cur[key] = 0
}

func putInTable(dataa string) string {
	data := strings.SplitN(dataa, "#", 3)
	program := data[0]
	procedures := strings.Split(data[1], "$")
	ip := data[2]

	for _, proc := range procedures {
		key := program + proc
		addEntry(key, ip)
		cur[key] = 0
	}

	return program
}

func handleConnection(conn net.Conn) {
	//p("Port mapper received connection")
	dec := gob.NewDecoder(conn)
	packet := pac{}
	dec.Decode(&packet)
	//fmt.Println("Received:", packet)
	switch packet.Ptype {
	case 0:
		fmt.Println("====================================================================\n[[[SERVER REGISTER RECEIVED]]]")
		ack := putInTable(packet.Content)
		//fmt.Println(packet.Content)
		fmt.Println(table)
		// acknowledge
		send1(conn, ack)
	case 2:
		// client request server location
		key := packet.Content
		//fmt.Println(packet.Content)
		// look up table
		_, ok := table[key]
		if ok {
			data := table[key][cur[key]]
			// shift current pointer
			cur[key] += 1
			cur[key] %= len(table[key])
			// send back
			send3(conn, data)
			fmt.Println("====================================================================\n[[[SENT REQUESTED LOCATION TO CLIENT]]]")

		} else {
			fmt.Println("Service not found in table!")
			return
		}

	default:
		// not for me, disgard
		fmt.Printf("here")
	}
	return
}

func main() {
	fmt.Println("====================================================================\n[[[PORT MAPPER STARTED]]]")
	fmt.Println(table)
	//create table in memory
	table = make(map[string][]string)
	cur = make(map[string]int)

	//p("Port mapper starting listening")
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()

	//p("Port mapper starting & writing its own location to file")
	file, _ := os.Create("../pmlocation") // write the file on afs
	port := getport(l)
	add := getmyip()
	location := add + ":" + port
	io.WriteString(file, location)

	for {
		conn, _ := l.Accept()     // this blocks until connection or error
		go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}

func p(str string) {
	fmt.Println(str)
}
