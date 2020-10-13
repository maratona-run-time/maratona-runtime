package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func Compile() {
	logFile, _ := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	log.SetOutput(logFile)
	jsonFile, openErr := os.Open("compilers.json")
	if openErr != nil {
		log.Fatal("Arquivo 'compilers.json' não pode ser aberto\n", openErr)
	}
	byteValueJSON, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
		log.Fatal("Arquivo 'compilers.json' não pode ser lido\n", readErr)
	}
	var result map[string]interface{}
	json.Unmarshal(byteValueJSON, &result)

	jsonFile.Close()
}

func main() {
	Compile()
}
