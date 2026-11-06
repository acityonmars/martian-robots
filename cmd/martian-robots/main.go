package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	engine "github.com/acityonmars/martian-robots/internal/engine"
	"github.com/acityonmars/martian-robots/internal/parser"
)

const maxLineBytes int = 1024
const maxLineNumber int = 1024

var debug = flag.Bool("debug", false, "enable debug output")

func cliUsage() {
	out := flag.CommandLine.Output()
	fmt.Fprintln(out, "Usage: martian-robots [-debug] [input-file]")
	fmt.Fprintln(out, "e.g. ./martian-robots -debug < sample.in")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Read from stdin when input-file is omitted. Press Ctrl+D on an empty line to finish.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Flags:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = cliUsage
	flag.Parse()
	logger := newLogger(*debug)

	// Read input name which could be "" for stdio, or a file name like ./testdata/sample.in.
	name, err := inputNameFromArgs(flag.Args())
	if err != nil {
		logAndExist(err)
	}
	reader, err := openInputIOReader(name, os.Stdin)
	if err != nil {
		logAndExist(err)
	}
	defer reader.Close()

	// Read raw input data into string slice.
	inputDataLines, err := readLines(reader)
	if err != nil {
		logAndExist(fmt.Errorf("read %s: %w", displayName(name), err))
	}

	// Parse input data, and create world and mission.
	world, mission, err := parser.ParseInputData(inputDataLines, logger)
	if err != nil {
		logAndExist(err)
	}

	// Create rover engine.
	marsEngine := engine.NewEngine(world, mission, logger)
	if err != nil {
		logAndExist(err)
	}

	// Run the engine.
	if err = marsEngine.Run(os.Stdout); err != nil {
		logAndExist(err)
	}
}

func displayName(name string) string {
	if name == "" {
		return "stdin"
	}
	return fmt.Sprintf("%q", name)
}

func inputNameFromArgs(args []string) (string, error) {
	switch len(args) {
	case 0:
		return "", nil
	case 1:
		return args[0], nil
	default:
		return "", fmt.Errorf("expected at most 1 input, got %d arguments", len(args))
	}
}

func openInputIOReader(name string, stdin io.Reader) (io.ReadCloser, error) {
	if name == "" {
		return io.NopCloser(stdin), nil
	}
	return os.Open(name)
}

func readLines(reader io.Reader) ([]string, error) {
	fileScanner := bufio.NewScanner(reader)
	fileScanner.Buffer(make([]byte, 0, maxLineBytes), maxLineBytes) // 1 KB per line
	var lines []string
	for lineNumber := 1; fileScanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(fileScanner.Text())
		if line == "" {
			return nil, fmt.Errorf("blank line in input at line %d is not accepted", lineNumber)
		}
		lines = append(lines, line)
		if len(lines) > maxLineNumber {
			return nil, fmt.Errorf("input exceeds %d lines", maxLineNumber)
		}
	}
	if err := fileScanner.Err(); err != nil {
		return nil, fmt.Errorf("input scanning failed: %w", err)
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
