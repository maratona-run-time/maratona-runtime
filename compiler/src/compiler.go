package compiler

import (
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

	logFile, _ := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	defer logFile.Close()
	log.SetOutput(logFile)
	commands := compilationCommand[compiler]
	_, execErr := exec.Command(commands[0], commands[1:]...).Output()
	if execErr != nil {
		log.Println("Erro na compilação\n", execErr)
		return "", execErr
	}
	return "program.out", nil
}
