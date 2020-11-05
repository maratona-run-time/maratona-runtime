package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

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

		status := make(chan []string)
		go executor.Execute("program.out", "tests/001.in", 2., status)
		res := <-status
		return string(res[0])
	})
	m.RunOnAddr(":8080")
}
