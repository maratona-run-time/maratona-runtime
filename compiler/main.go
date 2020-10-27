package main

import (
	"github.com/go-martini/martini"
	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	"github.com/martini-contrib/binding"

	"os"
)

type req struct {
	Program  string `json:"program"`
	Language string `json:"language"`
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.Json(req{}), func(r req) string {
		f, createErr := os.Create("program")
		defer f.Close()
		if createErr != nil {
			// log
			return "deu ruim"
		}
		f.WriteString(r.Program)
		ret, compilerErr := compiler.Compile(r.Language)
		if compilerErr != nil {
			// log
			return compilerErr.Error()
		}
		return ret
	})
	m.RunOnAddr(":8080")
}
