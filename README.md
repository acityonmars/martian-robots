# martian-robots

Martian robots simulation

## Usage

```sh
go build -o martian-robots ./cmd/martian-robots
./martian-robots [-debug] <input-file>
```

| Flag     | Default | Description         |
|----------|---------|---------------------|
| `-debug` | `false` | Enable debug output |

Flags must come before the input file (`./martian-robots -debug input.txt`). Go's `flag` package stops parsing at the
first non-flag argument, so `./martian-robots input.txt -debug` is rejected.

Run `./martian-robots -h` to print the usage.

## Assumptions

- Robots share a cell. (i.e. No colliding.)

## Design decisions

- Strict input data handling: If input data is invalid, fail immediately. Don't use partial data to prevent a potential
  misinterpretation of data.
- Execute all plans: However, don't stop executing plans if one plan fails. Instead, collect errors for failing plans
  and continue executing remaining plans.
- No metrics: Use clean logs over metrics. Consider using metrics if this app becomes a long-running application.
- No config file: Avoid using config file util this app requires something meaningful configuration. Use cleanly defined
  hardcoded values for now.
- Case-sensitive input data: Input data is case-sensitive.
- Treat no command as an error. 
