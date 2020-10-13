package main

import (
	"Maratona-Runtime/executor"
	"fmt"
	"os"
)

func main() {
	executable := "a.out"
	if len(os.Args[1:]) > 0 {
		executable = os.Args[1]
	}
	file := []string{fmt.Sprintf("./%s", executable)}

	// TODO: get input and output file names from command line
	status, err := executor.Verdict(file, "in", "out")
	fmt.Println(status)
	if err != nil {
		fmt.Println(err)
	}
}
