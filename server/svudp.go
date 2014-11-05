package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// +build dev
const SLICE = 5
const BURST = 2

// data packet sent to server
type DATA struct {
	packetNum    uint8        // which packet am i
	totalPackets uint8        // how many packets in total
	dataSlice    [SLICE]uint8 // the actual data
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

func byte2int(a []byte) [SLICE]uint8 {
	b := [SLICE]uint8{}
	for ind, val := range a {
		b[ind] = uint8(val)
	}
	return b
}

func parseargs(args []string) (rownum, colnum []int) {
	//args: [Math1 Min 1 1x5 1x1]
	//		[Math1 Mul 2 1x5 5x2 1x2]
	paranum, _ := strconv.Atoi(args[2])
	for i := 0; i < paranum; i++ {
		tmp := strings.Split(args[3+i], "x")
		t1, _ := strconv.Atoi(tmp[0])
		t2, _ := strconv.Atoi(tmp[1])
		rownum = append(rownum, t1)
		colnum = append(colnum, t2)
	}
	return rownum, colnum
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

func pointer2data(a []uint8) [SLICE]uint8 {
	arr := [SLICE]uint8{}
	for i := 0; i < len(a); i++ {
		arr[i] = a[i]
	}
	return arr
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

// send marshaled data to client
func senddata(con *net.UDPConn, addr *net.UDPAddr, packets []DATA) {
	for _, pac := range packets {
		// send data
		con.WriteToUDP(pac.concat(), addr)
		fmt.Println("Sent result:", pac)
	}
}

func handleUDP(con *net.UDPConn, args []string) {
	// receive all data
	buf := [512]byte{}
	ds := []DATA{}
	var n int
	var addr *net.UDPAddr
	for {
		n, addr, _ = con.ReadFromUDP(buf[0:])
		if n != 0 {
			d := DATA{packetNum: buf[0], totalPackets: buf[1], dataSlice: pointer2data(buf[2:n])}
			ds = append(ds, d)
			if len(ds)%BURST == 0 && len(ds) != 0 { // bulk ack
				con.WriteToUDP([]byte{buf[0]}, addr)
			}
			if uint8(len(ds)) == buf[1] { // got all data
				break
			}
		}
	}
	alldata := redo(ds)
	//demarshal using rownum, colnum
	rownum, colnum := parseargs(args)
	paras := [][][]byte{}
	tmp := 0
	for i := 0; i < len(rownum); i++ {
		para := demarshal(alldata[tmp:], rownum[i], colnum[i])
		tmp += rownum[i] * colnum[i]
		paras = append(paras, para)
	}
	//calculate&return
	switch args[1] {
	case "Min":
		fmt.Println("Demarshaled Data: ", paras[0])
		min := realMin(paras[0])
		fmt.Println("Min: ", min)
		con.WriteToUDP([]byte{1, 1, min}, addr)
		fmt.Println("[[[UDP CONNECTION CLOSED FOR " + args[1] + "]]]\n====================================================================")
	case "Max":
		fmt.Println("Demarshaled Data: ", paras[0])
		max := realMax(paras[0])
		fmt.Println("Max: ", max)
		con.WriteToUDP([]byte{1, 1, max}, addr)
		fmt.Println("[[[UDP CONNECTION CLOSED FOR " + args[1] + "]]]\n====================================================================")
	case "Multiply":
		fmt.Println("Demarshaled Data: ", paras)
		multiply := realMultiply(paras[0], paras[1], len(paras[0]), len(paras[0][0]), len(paras[1][0]))
		fmt.Println("Multiply: ", multiply)
		pacs := marshal(bytes.Join(multiply, []byte{}))
		senddata(con, addr, pacs)
		fmt.Println("[[[UDP CONNECTION CLOSED FOR " + args[1] + "]]]\n====================================================================")
	case "Sort":
		fmt.Println("Demarshaled Data: ", paras[0])
		sort := realSort(paras[0])
		fmt.Println("Sorted: ", sort)
		pacs := marshal(sort)
		senddata(con, addr, pacs)
		fmt.Println("[[[UDP CONNECTION CLOSED FOR " + args[1] + "]]]\n====================================================================")
	default:
	}
}

func realMin(para [][]byte) byte {
	min := byte(255)
	for i := 0; i < len(para); i++ {
		for j := 0; j < len(para[i]); j++ {
			if para[i][j] < min {
				min = para[i][j]
			}
		}
	}
	return min
}

func realMax(para [][]byte) byte {
	max := byte(0)
	for i := 0; i < len(para); i++ {
		for j := 0; j < len(para[i]); j++ {
			if para[i][j] > max {
				max = para[i][j]
			}
		}
	}
	return max
}

func realSort(para [][]uint8) []uint8 {
	sorted := make([]uint8, len(para[0]))
	copy(sorted, para[0])
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted); j++ {
			if sorted[i] < sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

func realMultiply(a, b [][]uint8, m, n, l int) [][]uint8 {
	c := [][]uint8{}
	var i, j, k int
	for i = 0; i < m; i++ {
		row := []uint8{}
		for k = 0; k < l; k++ {
			tmp := uint8(0)
			for j = 0; j < n; j++ {
				tmp += a[i][j] * b[j][k]
			}
			row = append(row, tmp)
		}
		c = append(c, row)
	}
	return c
}
