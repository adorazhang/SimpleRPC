// client
package main

import (
	"fmt"
	"os"
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

func main() {
	//args = os.Args
	//args = []string{"Math", "1", "Sort", "5"}
	args = []string{"Math", "1", "Multiply", "2", "3", "3", "1"}
	program = args[0]
	version = args[1]
	procedure = args[2]

	if args == nil || len(args) < 3 {
		Usage()
		return
	}
	if program != "Math" {
		p("Please use the only program: math for testing!")
		Usage()
		return
	}
	if version != "1" && version != "2" {
		p("We only support version 1 and 2 now!")
		Usage()
		return
	}

	//Client calling remote procedure as if local
	arr := []uint8{5, 2, 4, 16, 23, 6, 3, 9, 13, 6, 4, 123, 2, 6, 4, 13, 91, 16}
	m1 := [][]uint8{{1, 2, 3}, {3, 4, 5}}
	m2 := [][]uint8{{2}, {1}, {2}}
	switch procedure {
	case "Min":
		if len(args) != 4 {
			Usage()
			return
		}
		ret := Min(arr)
		fmt.Println("Min:", ret)
	case "Max":
		if len(args) != 4 {
			Usage()
			return
		}
		ret := Max(arr)
		fmt.Println("Max:", ret)
	case "Sort":
		if len(args) != 4 {
			Usage()
			return
		}
		ret := Sort(arr)
		fmt.Println("Sort: ", ret)
	case "Multiply":
		if len(args) != 7 {
			Usage()
			return
		}
		if args[4] != args[5] {
			p("Matrix dimension must match!")
		}
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
