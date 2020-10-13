package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Compile(compiler string) {
	logFile, _ := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	defer logFile.Close()
	log.SetOutput(logFile)
	jsonFile, openErr := os.Open("compilers.json")
	defer jsonFile.Close()
	if openErr != nil {
		log.Fatal("Arquivo 'compilers.json' não pode ser aberto\n", openErr)
	}
	byteValueJSON, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
		log.Fatal("Arquivo 'compilers.json' não pode ser lido\n", readErr)
	}
	var compilationCommand map[string][]string
	json.Unmarshal(byteValueJSON, &compilationCommand)

	commands := compilationCommand[compiler]
	_, execErr := exec.Command(commands[0], commands[1:]...).Output()
	if execErr != nil {
		log.Fatal("Erro compilacao\n", execErr)
	}
}

func main() {
	Compile("C")
}
