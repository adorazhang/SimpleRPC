// clientstub
package main

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
)

func portmapper() string {
	pm := conn2pm()
	loc := getserver(pm, program+version+procedure)
	pm.Close()
	return loc
}

func Min(arr []uint8) [][]uint8 {
	loc := portmapper()
	// done with the port mapper                                                                      //finished talking to port mapper
	ser := conn2sTCP(loc, program+version, procedure, "1", "1x"+strconv.Itoa(len(arr)), "1x1") //get udpconnection
	// done with TCP connection with the server
	fmt.Println("====================================================================\n[[[UDP CONNECTION ESTABLISHED]]]")
	packets := marshal(arr)
	senddata(ser, packets)
	results := getresult(ser, 1, 1)
	return results
}

func serialize(m1, m2 [][]uint8) []uint8 {
	res := []uint8{}
	for i := 0; i < len(m1); i++ {
		for j := 0; j < len(m1[i]); j++ {
			res = append(res, m1[i][j])
		}
	}
	for i := 0; i < len(m2); i++ {
		for j := 0; j < len(m2[i]); j++ {
			res = append(res, m2[i][j])
		}
	}
	return res
}

func Multiply(m1, m2 [][]uint8) [][]uint8 {
	m := strconv.Itoa(len(m1))
	n := strconv.Itoa(len(m1[0]))
	l := strconv.Itoa(len(m2[0]))
	loc := portmapper()
	// done with the port mapper                                                                      //finished talking to port mapper
	ser := conn2sTCP(loc, program+version, procedure, "2", m+"x"+n+"#"+n+"x"+l, m+"x"+l) //get udpconnection
	// done with TCP connection with the server
	fmt.Println("====================================================================\n[[[UDP CONNECTION ESTABLISHED]]]")
	arr := serialize(m1, m2)
	packets := marshal(arr)
	senddata(ser, packets)
	results := getresult(ser, len(m1), len(m2[0]))
	fmt.Println(results)
	return results
}

func Max(arr []uint8) [][]uint8 {
	loc := portmapper()
	// done with the port mapper                                                                      //finished talking to port mapper
	ser := conn2sTCP(loc, program+version, procedure, "1", "1x"+strconv.Itoa(len(arr)), "1x1") //get udpconnection
	// done with TCP connection with the server
	fmt.Println("====================================================================\n[[[UDP CONNECTION ESTABLISHED]]]")
	packets := marshal(arr)
	senddata(ser, packets)
	results := getresult(ser, 1, 1)
	return results
}

func Sort(arr []uint8) [][]uint8 {
	loc := portmapper()
	// done with the port mapper                                                                      //finished talking to port mapper
	ser := conn2sTCP(loc, program+version, procedure, "1", "1x"+strconv.Itoa(len(arr)), "1x"+strconv.Itoa(len(arr))) //get udpconnection
	// done with TCP connection with the server
	fmt.Println("====================================================================\n[[[UDP CONNECTION ESTABLISHED]]]")
	packets := marshal(arr)
	senddata(ser, packets)
	results := getresult(ser, 1, len(arr))
	return results
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

func getserver(conn net.Conn, key string) string {
	//send request packet
	send(conn, 2, key)
	//wait for response packet
	dec := gob.NewDecoder(conn)
	packet := PACKET{}
	dec.Decode(&packet)
	return packet.Content
}

func conn2sTCP(loc, progver, proce, numpara, paralen, returnlen string) *net.UDPConn {
	// Client stub connecting to server
	TCPconn, _ := net.Dial("tcp", loc)
	if TCPconn == nil {
		p("Cannot connect to server!")
		os.Exit(1)
	}
	// negotiate with the server
	data := progver + "#" + proce + "#" + numpara + "#" + paralen + "#" + returnlen // e.g. Math1#Multiply#2#2x3#3x1#2x1
	//fmt.Println(data)
	send(TCPconn, 4, data)

	//wait for UDP info
	dec := gob.NewDecoder(TCPconn)
	packet := PACKET{}
	dec.Decode(&packet)

	if packet.Ptype != 5 {
		p("Service not provided by that server")
		return nil
	}
	TCPconn.Close() // end talking in TCP
	fmt.Println("====================================================================\n[[[END TALKING IN TCP CONNECTION]]]")

	// got UDP info, establish connection
	laddr, _ := net.ResolveUDPAddr("udp", packet.Content)
	UDPconn, _ := net.DialUDP("udp", nil, laddr)

	return UDPconn
}
