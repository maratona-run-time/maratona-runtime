package main

import "github.com/go-martini/martini"

func createOrmServer() *martini.ClassicMartini {
	m := martini.Classic()
	setChallengeRoutes(m)
	return m
}

func main() {
	m := createOrmServer()
	m.RunOnAddr(":8080")
}
