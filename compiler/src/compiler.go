package compiler

import (
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os/exec"
)

func Compile(compiler string, fileName string, logger zerolog.Logger) (string, error) {
	compilationCommand := map[string][]string{
		"C":      {"gcc", fileName, "-o", "program.out"},
		"C++":    {"g++", fileName, "-o", "program.out"},
		"C++11":  {"g++", fileName, "-std=c++11", "-o", "program.out"},
		"Python": {"cp", fileName, "program.out"},
		"Go":     {"go", "build", "-o", "program.out", fileName},
	}

	shebangDict := map[string]string{
		"Python": "#!/usr/bin/env python3",
	}

	commands, compilerSupported := compilationCommand[compiler]
	if compilerSupported == false {
		logger.Info().
			Msg("Programming language not supported")
		return "", fmt.Errorf("Language '" + compiler + "' required for compilation not supported")
	}
	_, execErr := exec.Command(commands[0], commands[1:]...).Output()
	if execErr != nil {
		logger.Warn().
			Err(execErr).
			Msg("Compilation Error\n")

		return "", execErr
	}

	// Adiciona o Shebang
	if shebang, ok := shebangDict[compiler]; ok {
		code, readErr := ioutil.ReadFile("program.out")
		if readErr != nil {
			logger.Error().
				Err(readErr).
				Msg("An error happened while reading the file to add the Shebang\n")

			return "", readErr
		}
		executable := append([]byte(shebang+"\n"), code...)
		writeErr := ioutil.WriteFile("program.out", executable, 0755)
		if writeErr != nil {
			logger.Error().
				Err(writeErr).
				Msg("An error happened while writing in the file to add the Shebang\n")
			return "", writeErr
		}
	}

	logger.Debug().
		Msg("Compilation Finished")

	return "program.out", nil
}
