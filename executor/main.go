package main

import (
	"encoding/json"
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
	Inputs []*multipart.FileHeader `form:"inputs"`
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(f FileForm) []byte {
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

		os.Mkdir("inputs", 0700)

		for _, file := range f.Inputs {
			testFileName := fmt.Sprintf("inputs/%s", file.Filename)
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

		res := executor.Execute("program.out", "inputs", 2.)
		jsonResult, _ := json.Marshal(res)
		return jsonResult
	})
	m.RunOnAddr(":8080")
}
