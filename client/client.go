// client
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

var args []string
var program, version, procedure string

func Usage() {
	fmt.Println("USAGE: client.exe program_name version commands [arguments] ...")
	fmt.Println("The commands [arguments] are:")
	fmt.Println("Multiply\t<row1_length>\t<col1_length>\t<row2_length>\t<col2_length>")
	fmt.Println("Sort\t<array_length>")
	fmt.Println("Min\t<array_length>")
	fmt.Println("Max\t<array_length>")
	fmt.Println("Help")
}

func randomarr(len int) []uint32 {
	a := []uint32{}
	for i := 0; i < len; i++ {
		a = append(a, uint32(rand.Intn(100)))
	}
	return a
}

func randommatrix(row, col int) [][]uint32 {
	a := [][]uint32{}
	for i := 0; i < row; i++ {
		b := []uint32{}
		for i := 0; i < col; i++ {
			b = append(b, uint32(rand.Intn(100)))
		}
		a = append(a, b)
	}
	return a
}

func main() {
	args = os.Args
	if len(args) == 1 {
		Usage()
		return
	}
	program = args[1]
	version = args[2]
	procedure = args[3]
	if program != "Math" || program == "help" {
		p("Please use the only program: Math for testing!")
		Usage()
		return
	}
	if version != "1" && version != "2" {
		p("We only support version 1 and 2 now!")
		Usage()
		return
	}

	//args = []string{"Math", "1", "Sort", "36"}
	//args = []string{"Math", "1", "Multiply", "2", "3", "3", "1"}
	//Client calling remote procedure as if local
	switch procedure {
	case "Min":
		if len(args) != 5 {
			Usage()
			return
		}
		len, _ := strconv.Atoi(args[4])
		arr := randomarr(len)
		ret := Min(arr)
		fmt.Println("Min:", ret)
	case "Max":
		if len(args) != 5 {
			Usage()
			return
		}
		len, _ := strconv.Atoi(args[4])
		arr := randomarr(len)
		ret := Max(arr)
		fmt.Println("Max:", ret)
	case "Sort":
		if len(args) != 5 {
			Usage()
			return
		}
		len, _ := strconv.Atoi(args[4])
		arr := randomarr(len)
		ret := Sort(arr)
		fmt.Println("Sort: ", ret)
	case "Multiply":
		if len(args) != 8 {
			Usage()
			return
		}
		if args[5] != args[6] {
			p("Matrix dimension must match!")
		}
		m, _ := strconv.Atoi(args[4])
		n, _ := strconv.Atoi(args[5])
		l, _ := strconv.Atoi(args[7])
		m1 := randommatrix(m, n)
		m2 := randommatrix(n, l)
		ret := Multiply(m1, m2)
		fmt.Println("Multipy:")
		for i := 0; i < len(m1); i++ {
			fmt.Println(ret[i])
		}
	default:
		Usage()
		return
	}
}

func p(str string) {
	fmt.Println(str)
}

func printerr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
