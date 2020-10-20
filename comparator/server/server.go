package main

import (
	"github.com/go-martini/martini"
	"github.com/maratona-run-time/Maratona-Runtime/comparator"
	"github.com/martini-contrib/binding"
)

type req struct {
	S1 string `json:"s1"`
	S2 string `json:"s2"`
}

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
