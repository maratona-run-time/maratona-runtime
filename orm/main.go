package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Challenge struct {
	gorm.Model
	Title       string
	TimeLimit   int
	MemoryLimit int
}

func dbConnect() *gorm.DB {
	host := "localhost"
	port := "5432"
	user := "postgres"
	dbname := "mart"
	password := "password"

	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect database")
	}
	return db
}

func initialMigration() {
	db := dbConnect()
	db.AutoMigrate(&Challenge{})
}

func createChallenge(c *Challenge) {
	db := dbConnect()
	db.Create(c)
}

func readFirstChallenge(c *Challenge) {
	db := dbConnect()
	db.First(c)
}

func deleteChallenge(title string) {
	db := dbConnect()
	var c Challenge
	db.Delete(&c, "Title = ?", title)
}

func readAllChallenges(challenges *[]Challenge) {
	db := dbConnect()
	db.Find(challenges)
}

func test() {
	createChallenge(&Challenge{Title: "Teste"})
	deleteChallenge("Teste")
	var all []Challenge
	readAllChallenges(&all)
	for _, challenge := range all {
		fmt.Println(challenge.Title)
	}
}

func main() {
	initialMigration()
	test()
}