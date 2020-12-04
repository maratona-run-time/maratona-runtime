package orm

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"

	model "github.com/maratona-run-time/Maratona-Runtime/model"
)

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
		if err = db.AutoMigrate(&model.TestFile{}); err != nil {
			panic(err)
		}
		if err = db.AutoMigrate(&model.Challenge{}); err != nil {
			panic(err)
		}
	})
	return db
}

func CreateChallenge(challenge *model.Challenge) error {
	db := dbConnect()
	return db.Create(challenge).Error
}

func FindChallenge(id string) (model.Challenge, error) {
	db := dbConnect()
	var challenge model.Challenge
	err := db.Preload("Inputs").Preload("Outputs").First(&challenge, id).Error
	return challenge, err
}

func FindAllChallenges() ([]model.Challenge, error) {
	db := dbConnect()
	var challenges []model.Challenge
	err := db.Preload("Inputs").Preload("Outputs").Find(&challenges).Error
	return challenges, err
}
