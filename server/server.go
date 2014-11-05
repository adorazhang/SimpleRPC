// server
package main

import (
	"fmt"
	"os"
	"strings"
)

func Usage() {
	fmt.Println("USAGE: server.exe commands [arguments] ...")
	fmt.Println("The commands [arguments] are:")
	fmt.Println("Start\t<Program Name>\t<Version>")
	fmt.Println("Help")
}

func main() {
	//args := os.Args
	args := []string{"Start", "Math", "1"}
	choice := args[0]
	program := args[1]
	version := args[2]

	if args == nil || len(args) < 2 {
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

	switch choice {
	case "Start":
		if len(args) != 3 {
			Usage()
			return
		}
		Register(strings.Join(args[1:3], ""))
	default:
		Usage()
		return
	}
}

func p(str string) {
	fmt.Println(str)
}

func Min(arr []int) {

}

func Max() {
	fmt.Printf("")
}

func Sort() {
}

func Multiply() {
}

func printerr(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
