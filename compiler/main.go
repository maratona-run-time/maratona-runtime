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
			rs.WriteHeader(http.StatusBadRequest)
			rs.Write([]byte("An error occurred while trying to create a file named '" + fileName + "'"))
			return
		}
		program, pErr := req.Program.Open()
		if pErr != nil {
			rs.WriteHeader(http.StatusBadRequest)
			rs.Write([]byte("An error occurred while trying to open the received program"))
			return
		}
		io.Copy(f, program)
		f.Close()
		program.Close()
		ret, compilerErr := compiler.Compile(req.Language, fileName)
		if compilerErr != nil {
			rs.WriteHeader(http.StatusBadRequest)
			rs.Write([]byte("An error occurred while trying compile program in language '" + req.Language + "'"))
			return
		}
		http.ServeFile(rs, rq, ret)
	})
	m.RunOnAddr(":8080")
}
