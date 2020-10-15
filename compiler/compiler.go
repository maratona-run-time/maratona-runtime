package compiler

import (
	"log"
	"os"
	"os/exec"
)

func Compile(compiler string) (string, error) {
	compilationCommand := map[string][]string{
		"C":      {"gcc", "programa.c", "-o", "programa.out"},
		"C++":    {"g++", "programa.cpp", "-o", "programa.out"},
		"C++11":  {"g++", "programa.cpp", "-std=c++11", "-o", "programa.out"},
		"Python": {"cp", "programa.py", "programa.out"},
		"Go":     {"go", "build", "-o", "programa.out", "programa.go"},
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
	return "programa.out", nil
}
