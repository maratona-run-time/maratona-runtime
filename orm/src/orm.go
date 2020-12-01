package orm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Challenge struct {
	gorm.Model
	Title       string
	TimeLimit   int
	MemoryLimit int
}

var db *gorm.DB = nil
var once sync.Once

func dbConnect() *gorm.DB {
	once.Do(func() {
		host := "localhost"
		port := "5432"
		user := "postgres"
		dbname := "mart"
		password := "password"
		dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, dbname, password)
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
	})
	return db
}

func initialMigration() {
	db := dbConnect()
	db.AutoMigrate(&Challenge{})
}

func CreateChallenge(c *Challenge) {
	db := dbConnect()
	db.Create(c)
}

func FindChallenge(id string) Challenge {
	db := dbConnect()
	var c Challenge
	db.First(&c, id)
	return c
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

func Test() {
	CreateChallenge(&Challenge{Title: "Teste"})
	deleteChallenge("Teste")
	var all []Challenge
	readAllChallenges(&all)
	for _, challenge := range all {
		fmt.Println(challenge.Title)
	}
}
