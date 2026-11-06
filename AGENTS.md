# AGENTS.md

Take-home project: a stdin→stdout Go CLI simulating robots on a bounded Martian grid.
Quality bar: production-like, but do not over-engineer.

The original requirements are in [`docs/REQUIREMENT.md`](docs/REQUIREMENT.md). It is the source of truth: read it
before answering spec questions, and quote it when a rule's exact wording matters (e.g. "less than 100 characters").
If this file and the spec disagree, the spec wins; point out the mismatch.

## Priorities (limited time budget)

Focus on major issues and strong fundamentals; skip polishing.

- **Major issues** (raise first): wrong output or spec violations, crashes/panics, wrong exit codes, silent data loss
  (errors ignored), I/O boundary violations, missing or weak tests on core behavior.
- **Fundamentals** (worth the time): clear package boundaries, `io.Reader`/`io.Writer` at the edges, errors wrapped with
  context, golden-file and edge-case tests, documented assumptions and scale limits.
- **Polishing** (skip unless asked or trivial in passing): naming, comment wording, message phrasing, import aliases,
  formatting beyond `gofmt`, micro-optimizations.
- In reviews, list major issues only. Mention minor ones in a single line at most, or not at all.

## Working style

- Be concise.
- Give brief reasoning with each suggestion.
- Avoid assumptions; when something is unclear, ask.
- Ask for confirmation before making changes.

## Commands

```sh
go build ./...
go vet ./...
gofmt -l .                                   # must print nothing
go test ./...
go run ./cmd/martian-robots < testdata/sample.in
```

## Layout

- `cmd/martian-robots/` — entrypoint only: wires stdin/stdout/exit code, no logic.
- `internal/` — domain code. Prefer a named package (e.g. `internal/sim`) over files directly in `internal`.
- `testdata/` — input/expected-output fixtures (golden files).

## Domain rules (summary of `docs/REQUIREMENT.md`)

- Grid is `0,0` to the given upper-right corner; max coordinate is 50; instruction strings are < 100 chars.
- Orientations N/S/E/W; North is `y+1`. Commands: `L`, `R`, `F`.
- A robot moving off the grid is `LOST`; it leaves a scent at its last on-grid position.
- A scent makes later robots ignore any `F` that would leave the grid from that same point (scent is keyed by position only).
- Robots run sequentially; output one line per robot: `x y O` plus ` LOST` if lost.
- New command types must be addable without changing the robot execution loop (command interface or registry keyed by letter).

## Conventions

- Keep the simplest solution that works; document scale limits in the README instead of building for them.
- Explicit types for all functions and domain values (`Orientation`, `Point`, not raw strings/ints).
- Core logic takes `io.Reader`/`io.Writer`; only `main` touches `os.Stdin`, `os.Stdout`, `os.Exit`.
- Wrap errors with context (`fmt.Errorf("line %d: %w", n, err)`); errors go to stderr with a non-zero exit.
- Tests use real inputs and real implementations. No mocking libraries.
- Standard library only unless there is a clear reason.

## Do not

- Commit, or run destructive commands, without asking the user first.
- Add Docker, config files, logging/metrics frameworks, or CLI flag libraries.
- Commit build artifacts (`martian-robots` binary) or `.idea/`.