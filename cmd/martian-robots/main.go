package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/acityonmars/martian-robots/internal/display"
	engine "github.com/acityonmars/martian-robots/internal/engine"
	"github.com/acityonmars/martian-robots/internal/parser"
)

const maxLineBytes int = 1024
const maxLineNumber int = 1024
const maxNumberOfRobots = 100

func main() {

	// get the file name from argument
	fileName, err := getInputFileNameFromArgs()
	if err != nil {
		logAndExist(err)
	}

	// read the file into lines
	inputData, err := readFileToLines(fileName)
	if err != nil {
		logAndExist(err)
	}

	// parse the string
	worldUpperRight, mission, err := parser.ParseMissionData(inputData)
	if err != nil {
		logAndExist(err)
	}

	worldDisplay := display.NewConsoleDisplay()
	if err = worldDisplay.Init(*worldUpperRight); err != nil {
		logAndExist(err)
	}

	worldEngine := engine.NewEngine()
	if err = worldEngine.Init(worldDisplay); err != nil {
		logAndExist(err)
	}

	if err = worldEngine.Run(*mission); err != nil {
		logAndExist(err)
	}
}

func getInputFileNameFromArgs() (string, error) {
	args := os.Args

	argsCnt := len(args)
	if argsCnt != 2 {
		return "", errors.New("incorrect number of arguments provided")
	}

	return args[1], nil
}

func readFileToLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Buffer(make([]byte, 0, maxLineBytes), maxLineNumber) // approx. 10MB
	var lines []string
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if len(lines) == maxLineNumber {
			return nil, fmt.Errorf("file %q exceeds %d lines", fileName, maxLineNumber)
		}
	}
	if err = fileScanner.Err(); err != nil {
		return nil, fmt.Errorf("reading file %q failed: %w", fileName, err)
	}

	return lines, nil
}

func logAndExist(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
