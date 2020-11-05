package main

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/go-martini/martini"
	executor "github.com/maratona-run-time/Maratona-Runtime/executor/src"
	"github.com/martini-contrib/binding"
)

// FileForm define o tipo de dados esperado no POST.
// Recebe um arquivo binário e um conjunto de arquivos de entrada.
type FileForm struct {
	Binary *multipart.FileHeader   `form:"binary"`
	Tests  []*multipart.FileHeader `form:"tests"`
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(f FileForm) string {
		receivedFile, rErr := f.Binary.Open()
		if rErr != nil {
			panic(rErr)
		}

		binaryFile, bErr := os.Create("program.out")
		if bErr != nil {
			panic(bErr)
		}

		exeErr := os.Chmod("program.out", 0777)
		if exeErr != nil {
			panic(exeErr)
		}

		io.Copy(binaryFile, receivedFile)

		binaryFile.Close()
		receivedFile.Close()

		os.Mkdir("tests", 0700)
		for i, file := range f.Tests {
			testFileName := fmt.Sprintf("tests/%03d.in", i+1)
			testFile, testFileErr := os.Create(testFileName)
			if testFileErr != nil {
				panic(testFileErr)
			}
			defer testFile.Close()
			receivedTestFile, rfErr := file.Open()
			if rfErr != nil {
				panic(rfErr)
			}
			defer receivedTestFile.Close()
			io.Copy(testFile, receivedTestFile)
		}

		ch := make(chan []byte)
		chErr := make(chan error)
		ctx, cancel := context.WithTimeout(context.Background(), (time.Second * 2))
		defer cancel()
		go executor.Execute(ctx, "program.out", "tests/001.in", ch, chErr)
		select {
		case res := <-ch:
			return string(res)
		case err := <-chErr:
			panic(err)
		}
	})
	m.RunOnAddr(":8080")
}
