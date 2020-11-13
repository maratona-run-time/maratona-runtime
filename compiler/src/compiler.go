package compiler

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func Compile(compiler string) (string, error) {
	compilationCommand := map[string][]string{
		"C":      {"gcc", "program.c", "-o", "program.out"},
		"C++":    {"g++", "program.cpp", "-o", "program.out"},
		"C++11":  {"g++", "program.cpp", "-std=c++11", "-o", "program.out"},
		"Python": {"cp", "program.py", "program.out"},
		"Go":     {"go", "build", "-o", "program.out", "program.go"},
	}

	shebangDict := map[string]string{
		"Python": "#!/usr/bin/env python3",
	}

	logFile, _ := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	defer logFile.Close()
	log.SetOutput(logFile)
	commands := compilationCommand[compiler]
	_, execErr := exec.Command(commands[0], commands[1:]...).Output()
	if execErr != nil {
		log.Println("Erro na compilação\n", execErr)
		return "", execErr
	}

	// Adiciona o Shebang
	if shebang, ok := shebangDict[compiler]; ok {
		code, readErr := ioutil.ReadFile("program.out")
		if readErr != nil {
			log.Fatalln("Erro durante a leitura do arquivo na hora de adicionar Shebang\n", readErr)
		}
		executable := append([]byte(shebang+"\n"), code...)
		writeErr := ioutil.WriteFile("program.out", executable, 0755)
		if writeErr != nil {
			log.Fatalln("Erro durante a escrita do arquivo na hora de adicionar Shebang\n", writeErr)
		}
	}

	return "program.out", nil
}
