package main

import (
	"bytes"
	"strings"
	"testing"
)

// sampleInput and sampleOutput are the sample from REQUIREMENT.md.
const sampleInput string = `5 3
1 1 E
RFRFRFRF
3 2 N
FRRFLLFFRRFLL
0 3 W
LLFFFLFLFL
`

const sampleOutput string = `1 1 E
3 3 N LOST
2 3 S
`

// result holds everything run() produced, so each test can assert on what matters to it.
type result struct {
	code   int
	stdout string
	stderr string
}

// runCLI calls run() in-process with real readers/writers.
func runCLI(t *testing.T, args []string, stdin string) result {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := run(args, strings.NewReader(stdin), &stdout, &stderr)
	return result{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func Test_Running_CorrectInput_Results_Expected_Output(t *testing.T) {
	t.Run("Test sample from stdin", func(t *testing.T) {
		got := runCLI(t, nil, sampleInput)

		if got.code != exitOK {
			t.Fatalf("exit code = %d, want %d; stderr: %s", got.code, exitOK, got.stderr)
		}
		if got.stdout != sampleOutput {
			t.Errorf("stdout:\n%s\nwant:\n%s", got.stdout, sampleOutput)
		}
	})

	t.Run("Test sample from file", func(t *testing.T) {
		got := runCLI(t, []string{"../../testdata/sample.in"}, "")

		if got.code != exitOK {
			t.Fatalf("exit code = %d, want %d; stderr: %s", got.code, exitOK, got.stderr)
		}
		if got.stdout != sampleOutput {
			t.Errorf("stdout:\n%s\nwant:\n%s", got.stdout, sampleOutput)
		}
	})
}

func Test_IncorrectArgs_Fails(t *testing.T) {
	got := runCLI(t, []string{"a.in", "b.in"}, "")

	if got.code == exitOK {
		t.Errorf("exit code = %d, want non-zero", got.code)
	}
	if got.stdout != "" {
		t.Errorf("stdout = %q, want empty", got.stdout)
	}
}
