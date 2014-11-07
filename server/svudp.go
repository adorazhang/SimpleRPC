package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
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
func pointer2data(a []byte) [SLICE]uint32 {
	arr := [SLICE]uint32{}
	ind := 0
	for i := 0; i < SLICE; i++ {
		arr[i] = binary.LittleEndian.Uint32(a[ind : ind+4])
		ind += 4
	}
	return arr
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
	buf := [BUFFERSIZE]byte{}
	ds := []DATA{}
	var n int
	var addr *net.UDPAddr
	for {
		n, addr, _ = con.ReadFromUDP(buf[0:])
		if n != 0 {
			fmt.Println(binary.LittleEndian.Uint32(buf[0:4]), binary.LittleEndian.Uint32(buf[4:8]), pointer2data(buf[8:n]))
			d := DATA{packetNum: binary.LittleEndian.Uint32(buf[0:4]), totalPackets: binary.LittleEndian.Uint32(buf[4:8]), dataSlice: pointer2data(buf[8:n])}
			ds = append(ds, d)
			if len(ds)%BURST == 0 && len(ds) != 0 { // bulk ack
				con.WriteToUDP(buf[0:4], addr)
			}
			if uint32(len(ds)) == d.totalPackets { // got all data
				break
			}
		}
	}
	alldata := redo(ds)
	//demarshal using rownum, colnum
	rownum, colnum := parseargs(args)
	paras := [][][]uint32{}
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
		pacs := marshal(min)
		senddata(con, addr, pacs)
		fmt.Println("[[[UDP CONNECTION CLOSED FOR " + args[1] + "]]]\n====================================================================")
	case "Max":
		fmt.Println("Demarshaled Data: ", paras[0])
		max := realMax(paras[0])
		fmt.Println("Max: ", max)
		pacs := marshal(max)
		senddata(con, addr, pacs)
		fmt.Println("[[[UDP CONNECTION CLOSED FOR " + args[1] + "]]]\n====================================================================")
	case "Multiply":
		fmt.Println("Demarshaled Data: ", paras)
		multiply := realMultiply(paras[0], paras[1], len(paras[0]), len(paras[0][0]), len(paras[1][0]))
		fmt.Println("Multiply: ", multiply)
		pacs := marshal(multiply)
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

func realMin(para [][]uint32) []uint32 {
	min := uint32(9999999)
	for i := 0; i < len(para); i++ {
		for j := 0; j < len(para[i]); j++ {
			if para[i][j] < min {
				min = para[i][j]
			}
		}
	}
	return []uint32{min}
}

func realMax(para [][]uint32) []uint32 {
	max := uint32(0)
	for i := 0; i < len(para); i++ {
		for j := 0; j < len(para[i]); j++ {
			if para[i][j] > max {
				max = para[i][j]
			}
		}
	}
	return []uint32{max}
}

func realSort(para [][]uint32) []uint32 {
	sorted := make([]uint32, len(para[0]))
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

func realMultiply(a, b [][]uint32, m, n, l int) []uint32 {
	c := []uint32{}
	var i, j, k int
	for i = 0; i < m; i++ {
		for k = 0; k < l; k++ {
			tmp := uint32(0)
			for j = 0; j < n; j++ {
				tmp += a[i][j] * b[j][k]
			}
			c = append(c, tmp)
		}
	}
	return c
}
