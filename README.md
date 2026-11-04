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

- Input data handling: If input data is invalid, fail early. Don't use partial data to prevent a potential
  misinterpretation of data. (i.e. all or nothing)
- However, don't stop executing plans if one plan fails. Instead, collect errors for failing plans and continue
  executing remaining plans.
- Decouple "display" from "engine". Make "engine" re-usable, and display method flexible.
- No metrics. Use clean logs over metrics. Consider using metrics if this app becomes a long-running application.
- No config file. Avoid using config file util this app requires something meaningful configuration. Use cleanly defined
  hardcoded values for now.
- Input data is case-sensitive. This makes parsing simpler.
- No command is an error. This makes command handling simpler.
