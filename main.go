package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	actualOutput := make(chan []byte)
	errorOutput := make(chan error)

	ctx, errContext := timerContext()

	if errContext != nil {
		fmt.Println(errContext)
		return
	}

	executable := "a.out"
	if len(os.Args[1:]) > 0 {
		executable = os.Args[1]
	}
	file := []string{fmt.Sprintf("./%s", executable)}

	inputFile, errInFile := os.Open("in")

	if errInFile != nil {
		fmt.Println(errInFile)
		return
	}

	go execute(ctx, file, inputFile, actualOutput, errorOutput)
	select {
	case <-ctx.Done():
		fmt.Println("deu tle")
		return
	case err := <-errorOutput:
		fmt.Println("deu rte")
		fmt.Println("%s", err)
		return
	case out := <-actualOutput:
		fmt.Println("Compara as saidas")
		expectedData, _ := ioutil.ReadFile("out")
		expectedOut := string(expectedData)
		programOut := string(out)
		if strings.EqualFold(programOut, expectedOut) {
			fmt.Println("deu ac")
		} else {
			fmt.Println("deu wa")
		}
		return
	}
}

func execute(ctx context.Context, executable []string, inputFile *os.File, output chan<- []byte, errorOutput chan<- error) {
	cmd := exec.CommandContext(ctx, executable[0], executable[1:]...)
	cmd.Stdin = inputFile
	fmt.Println("Pegando o output..")
	programOutput, err := cmd.Output() // Nao ta conseguindo pegar o output de um arquivo que tenha dado RTE, tirar o case da linha 23 resolve
	fmt.Println("pegou")
	fmt.Println(err)
	fmt.Println(programOutput)
	if err != nil {
		errorOutput <- err
		return
	}
	output <- programOutput
}

func timerContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
}
