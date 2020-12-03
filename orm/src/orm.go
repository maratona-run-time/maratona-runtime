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
	Inputs      []Input `gorm:"ForeignKey:ChallengeID"`
}

type Input struct {
	gorm.Model
	Filename    string
	Content     []byte
	ChallengeID uint
}

var db *gorm.DB = nil
var once sync.Once

func dbConnect() *gorm.DB {
	once.Do(func() {
		host := "postgres"
		port := "5432"
		user := "postgres"
		dbname := "mart"
		password := "password"
		dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", host, port, user, dbname, password)
		var err error
		if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&Input{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&Challenge{}); err != nil {
			panic(err)
		}
	})
	return db
}

func CreateChallenge(challenge *Challenge) error {
	db := dbConnect()
	return db.Create(challenge).Error
}

func FindChallenge(id string) (Challenge, error) {
	db := dbConnect()
	var challenge Challenge
	err := db.Preload("Inputs").First(&challenge, id).Error
	return challenge, err
}

func FindAllChallenges() ([]Challenge, error) {
	db := dbConnect()
	var challenges []Challenge
	err := db.Preload("Inputs").Find(&challenges).Error
	return challenges, err
}
