package main

import (
	"net"
	//	"strconv"
	"bytes"
	"fmt"
)

// +build dev
const SLICE = 5
const BURST = 2

// data packet sent to server
type DATA struct {
	packetNum    byte        // which packet am i
	totalPackets byte        // how many packets in total
	dataSlice    [SLICE]byte // the actual data
}

func (a DATA) concat() []byte {
	b := [2 + SLICE]byte{}
	b[0] = a.packetNum
	b[1] = a.totalPackets
	for i := 0; i < SLICE; i++ {
		b[i+2] = a.dataSlice[i]
	}
	return b[:]
}

// fill up the packet
func fillup(a []uint8) [SLICE]byte {
	arr := [SLICE]byte{}
	for i := 0; i < len(a); i++ {
		arr[i] = a[i]
	}
	return arr
}

// split an array into data packets
func marshal(a []uint8) []DATA {
	totalPackets := byte(len(a)/SLICE) + 1
	ds := []DATA{}
	d := DATA{}
	for i := 0; i < len(a); i += SLICE {
		if i+SLICE > len(a) { // the last packet
			d = DATA{byte(i/SLICE) + 1, totalPackets, fillup(a[i:])}
		} else {
			d = DATA{byte(i/SLICE) + 1, totalPackets, fillup(a[i : i+SLICE])}
		}
		ds = append(ds, d)
	}
	return ds
}

// send marshaled data to server
func senddata(con *net.UDPConn, packets []DATA) {
	buffer := make([]byte, 512)
	for _, pac := range packets {
		// send data
		fmt.Println("Sent:", pac)
		con.Write(pac.concat())
		if pac.packetNum%BURST == 0 && pac.packetNum != 0 {
			// wait for ack
			n, _ := con.Read(buffer[0:])
			if n > 0 {
				fmt.Println("ACK:", buffer[0])
				if buffer[0] == pac.totalPackets {
					return
				}
				if buffer[0] == pac.packetNum {
					continue
				} else {
					fmt.Println("Packet Lost!")
					break
				}
			}
		}
	}
}

func redo(data []DATA) []byte {
	a := [][]byte{}
	for i := uint8(1); i <= uint8(len(data)); i++ {
		for _, val := range data {
			if i == val.packetNum {
				a = append(a, val.dataSlice[0:SLICE])
				break
			}
		}
	}
	return bytes.Join(a, []byte{})
}

func pointer2data(a []byte) [SLICE]byte {
	arr := [SLICE]byte{}
	for i := 0; i < len(a); i++ {
		arr[i] = a[i]
	}
	return arr
}

func demarshal(arr []uint8, row, col int) [][]uint8 {
	b := [][]uint8{}
	for i := 0; i < row; i++ {
		tmp := make([]uint8, col)
		for j := 0; j < col; j++ {
			tmp[j] = arr[i*col+j]
		}
		b = append(b, tmp)
	}
	return b
}

func getresult(con *net.UDPConn, row, col int) [][]uint8 {
	// receive all data
	buf := [512]byte{}
	ds := []DATA{}
	for {
		n, _ := con.Read(buf[0:])
		if n != 0 {
			d := DATA{packetNum: buf[0], totalPackets: buf[1], dataSlice: pointer2data(buf[2:n])}
			ds = append(ds, d)
			fmt.Println("Received result:", d)
			if uint8(len(ds)) == buf[1] { // got all data
				break
			}
		}
	}
	con.Close() //finished RPC
	alldata := redo(ds)
	res := [][]uint8{}
	for i := 0; i < row; i++ {
		res = demarshal(alldata, row, col)
	}
	return res
}
