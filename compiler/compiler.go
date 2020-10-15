package compiler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Compile(compiler string) int {
	logFile, _ := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	defer logFile.Close()
	log.SetOutput(logFile)
	jsonFile, openErr := os.Open("compilers.json")
	defer jsonFile.Close()
	if openErr != nil {
		log.Println("Arquivo 'compilers.json' não pode ser aberto\n", openErr)
		return 1
	}
	byteValueJSON, readErr := ioutil.ReadAll(jsonFile)
	if readErr != nil {
		log.Println("Arquivo 'compilers.json' não pode ser lido\n", readErr)
		return 1
	}
	var compilationCommand map[string][]string
	json.Unmarshal(byteValueJSON, &compilationCommand)

	commands := compilationCommand[compiler]
	_, execErr := exec.Command(commands[0], commands[1:]...).Output()
	if execErr != nil {
		log.Println("Erro na compilação\n", execErr)
		return 1
	}
	return 0
}
