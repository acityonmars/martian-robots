package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	engine "github.com/acityonmars/martian-robots/internal/engine"
	"github.com/acityonmars/martian-robots/internal/parser"
)

const maxLineBytes int = 1024
const maxLineNumber int = 1024

var debug = flag.Bool("debug", false, "enable debug output")

func main() {
	flag.Parse()
	logger := newLogger(*debug)

	// Read the input file path from arguments.
	fileName, err := getInputFileNameFromArgs()
	if err != nil {
		logAndExist(err)
	}

	// Read raw input data.
	inputData, err := readFileToLines(fileName)
	if err != nil {
		logAndExist(err)
	}

	// Parse input data.
	world, mission, err := parser.ParseInputData(inputData, logger)
	if err != nil {
		logAndExist(err)
	}

	// Create rover engine.
	marsEngine := engine.NewEngine(world, mission, logger)
	if err != nil {
		logAndExist(err)
	}

	// Run the engine.
	if err = marsEngine.Run(); err != nil {
		logAndExist(err)
	}
}

func getInputFileNameFromArgs() (string, error) {
	args := flag.Args()

	argsCnt := len(args)
	if argsCnt != 1 {
		return "", errors.New("incorrect number of arguments provided")
	}

	return args[0], nil
}

func readFileToLines(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Buffer(make([]byte, 0, maxLineBytes), maxLineBytes) // approx. 10MB
	var lines []string
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)
		if len(lines) > maxLineNumber {
			return nil, fmt.Errorf("file %q exceeds %d lines", fileName, maxLineNumber)
		}
	}
	if err = fileScanner.Err(); err != nil {
		return nil, fmt.Errorf("reading file %q failed: %w", fileName, err)
	}

	return lines, nil
}

func newLogger(debug bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

func logAndExist(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
