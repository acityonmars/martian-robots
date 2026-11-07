package main

import (
	"bufio"
	"errors"
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

const (
	exitOK    = 0
	exitError = 1
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	// NOTE: Why local flag set? Each run() call gets fresh flags so tests don't leak state or panic on redefinition.
	// NOTE: "flag.ContinueOnError": Parse returns errors instead of os.Exit. This is to make run() testable.
	fs := flag.NewFlagSet("martian-robots", flag.ContinueOnError)
	fs.SetOutput(stderr)
	debug := fs.Bool("debug", false, "enable debug, verbose logging")
	fs.Usage = func() { cliUsage(fs) }

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK
		}
		return exitError
	}
	logger := newLogger(*debug, stderr)

	// Read input name which could be "" for stdio, or a file name like ./testdata/sample.in.
	name, err := inputNameFromArgs(fs.Args())
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitError
	}
	reader, err := openInputIOReader(name, stdin)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitError
	}
	defer reader.Close()

	// Read raw input data into string slice.
	inputDataLines, err := readLines(reader)
	if err != nil {
		fmt.Fprintln(stderr, fmt.Errorf("read %s: %w", displayName(name), err))
		return exitError
	}

	if *debug {
		printInputData(stderr, inputDataLines)
	}

	// Parse input data, and create world and robot plans.
	world, robotPlans, err := parser.ParseInputData(inputDataLines, logger)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitError
	}

	// Create rover engine.
	marsEngine := engine.NewEngine(world, robotPlans, logger)

	// Run the engine.
	if err = marsEngine.Run(stdout); err != nil {
		fmt.Fprintln(stderr, err)
		return exitError
	}

	return exitOK
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

func newLogger(debug bool, stderr io.Writer) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(stderr, &slog.HandlerOptions{Level: level}))
}

// printInputData presents raw input for debugging.
func printInputData(stderr io.Writer, rawInput []string) {
	fmt.Fprintln(stderr, "\n-------")
	for _, line := range rawInput {
		fmt.Fprintln(stderr, line)
	}
	fmt.Fprintln(stderr, "-------\n")
}

func cliUsage(fs *flag.FlagSet) {
	out := fs.Output()
	fmt.Fprintln(out, "Usage: martian-robots [-debug] [input-file]")
	fmt.Fprintln(out, "e.g. ./martian-robots -debug < sample.in")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Read from stdin when input-file is omitted. Press Ctrl+D on an empty line to finish.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Flags:")
	fs.PrintDefaults()
}
