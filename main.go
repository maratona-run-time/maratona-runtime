package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	actualOutput := make(chan []byte)
	errorOutput := make(chan error)

	ctx, _ := timerContext()

	executable := "a.out"
	if len(os.Args[1:]) > 0 {
		executable = os.Args[1]
	}
	file := []string{fmt.Sprintf("./%s", executable)}

	inputFile, _ := os.Open("in")

	go execute(ctx, file, inputFile, actualOutput, errorOutput)
	select {
	case <-ctx.Done(): // Tirando esse case os programas de RTE sao avaliados corretamente
		fmt.Println("deu tle")
		return
	case err := <-errorOutput:
		fmt.Println("deu rte")
		fmt.Println("%s", err)
		return
	case out := <-actualOutput:
		fmt.Println("Compara as saidas")
		stringOutput := string(out)
		fmt.Println(stringOutput)
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
	return context.WithDeadline(context.Background(), time.Now().Add(2000000))
}
