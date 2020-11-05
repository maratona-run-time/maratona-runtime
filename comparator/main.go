package main

import (
	"github.com/go-martini/martini"
	comparator "github.com/maratona-run-time/Maratona-Runtime/comparator/src"
	"github.com/martini-contrib/binding"
)

type req struct {
	S1 string `json:"s1"`
	S2 string `json:"s2"`
}

// POST /
// 	Espera parâmetros JSON 's1' e 's2' com as duas strings que devem ser comparadas.
//	Retorna "AC" se ambas strings são iguais, "WA" caso sejam diferentes.
func main() {
	m := martini.Classic()
	m.Post("/", binding.Json(req{}), func(r req) string {
		if comparator.Compare(r.S1, r.S2) {
			return "AC"
		}
		return "WA"
	})
	m.RunOnAddr(":8080")
}
