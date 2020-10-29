package main

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-martini/martini"
	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	"github.com/martini-contrib/binding"

	"os"
)

type FileForm struct {
	Language string                `form:"language"`
	Program  *multipart.FileHeader `form:"program"`
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, req FileForm) {
		fileName := req.Program.Filename
		f, createErr := os.Create(fileName)
		if createErr != nil {
			panic(createErr)
		}
		program, pErr := req.Program.Open()
		if pErr != nil {
			panic(pErr)
		}
		io.Copy(f, program)
		f.Close()
		program.Close()
		ret, compilerErr := compiler.Compile(req.Language)
		if compilerErr != nil {
			panic(compilerErr)
		}
		http.ServeFile(rs, rq, ret)
	})
	m.RunOnAddr(":8080")
}
