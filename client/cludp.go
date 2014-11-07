package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
)

// +build dev
const SLICE = 5
const BURST = 2
const BUFFERSIZE = 1024

// data packet sent to server
type DATA struct {
	packetNum    uint32        // which packet am i
	totalPackets uint32        // how many packets in total
	dataSlice    [SLICE]uint32 // the actual data
}

func (a DATA) concat() []byte {
	tmp := [][]byte{}

	t1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(t1, uint32(a.packetNum))
	tmp = append(tmp, t1)

	t2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(t2, uint32(a.totalPackets))
	tmp = append(tmp, t2)

	for i := 0; i < SLICE; i++ {
		t3 := make([]byte, 4)
		binary.LittleEndian.PutUint32(t3, uint32(a.dataSlice[i]))
		tmp = append(tmp, t3)
	}
	return bytes.Join(tmp, []byte{})
}

// fill up the packet
func fillup(a []uint32) [SLICE]uint32 {
	arr := [SLICE]uint32{}
	for i := 0; i < len(a); i++ {
		arr[i] = a[i]
	}
	return arr
}

// split an array into data packets
func marshal(a []uint32) []DATA {
	totalPackets := uint32(math.Ceil(float64(len(a)) / SLICE))
	ds := []DATA{}
	d := DATA{}
	for i := 0; i < len(a); i += SLICE {
		if i+SLICE > len(a) { // the last packet
			d = DATA{uint32(i/SLICE + 1), uint32(totalPackets), fillup(a[i:])}
		} else {
			d = DATA{uint32(i/SLICE + 1), uint32(totalPackets), fillup(a[i : i+SLICE])}
		}
		ds = append(ds, d)
	}
	return ds
}

// send marshaled data to server
func senddata(con *net.UDPConn, packets []DATA) {
	buffer := make([]byte, BUFFERSIZE)
	for _, pac := range packets {
		// send data
		fmt.Println("Sent:", pac)
		con.Write(pac.concat())
		if pac.packetNum%BURST == 0 && pac.packetNum != 0 {
			// wait for ack
			n, _ := con.Read(buffer[0:])
			if n > 0 {
				fmt.Println("ACK:", binary.LittleEndian.Uint32(buffer[0:4]))
				if binary.LittleEndian.Uint32(buffer[0:4]) == uint32(pac.totalPackets) {
					return
				}
				if binary.LittleEndian.Uint32(buffer[0:4]) == uint32(pac.packetNum) {
					continue
				} else {
					fmt.Println("Packet Lost!")
					break
				}
			}
		}
	}
}

func redo(data []DATA) []uint32 {
	a := []uint32{}
	for i := uint32(1); i <= uint32(len(data)); i++ {
		for _, val := range data {
			if i == val.packetNum {
				for j := 0; j < SLICE; j++ {
					a = append(a, val.dataSlice[j])
				}
				break
			}
		}
	}
	return a
}
func pointer2data(a []byte) [SLICE]uint32 {
	arr := [SLICE]uint32{}
	ind := 0
	for i := 0; i < SLICE; i++ {
		arr[i] = binary.LittleEndian.Uint32(a[ind : ind+4])
		ind += 4
	}
	return arr
}

func demarshal(arr []uint32, row, col int) [][]uint32 {
	b := [][]uint32{}
	for i := 0; i < row; i++ {
		tmp := make([]uint32, col)
		for j := 0; j < col; j++ {
			tmp[j] = arr[i*col+j]
		}
		b = append(b, tmp)
	}
	return b
}

func getresult(con *net.UDPConn, row, col int) [][]uint32 {
	// receive all data
	buf := [BUFFERSIZE]byte{}
	ds := []DATA{}
	for {
		n, _ := con.Read(buf[0:])
		if n != 0 {
			d := DATA{packetNum: binary.LittleEndian.Uint32(buf[0:4]), totalPackets: binary.LittleEndian.Uint32(buf[4:8]), dataSlice: pointer2data(buf[8:n])}
			ds = append(ds, d)
			fmt.Println("Received result:", d)
			if uint32(len(ds)) == binary.LittleEndian.Uint32(buf[4:8]) { // got all data
				break
			}
		}
	}
	con.Close() //finished RPC
	alldata := redo(ds)
	res := [][]uint32{}
	for i := 0; i < row; i++ {
		res = demarshal(alldata, row, col)
	}
	return res
}
