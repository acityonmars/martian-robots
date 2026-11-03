# martian-robots
Martian robots simulation


## Assumptions
- If input data is invalid, don't use partial data to prevent  a potential misinterpretation of data. (i.e. all or nothing)
- This application may grow with additional features.
- We may want to use a different display method other than stdout in the future.

## Design decisions
- Decouple "display" from "engine". Make "engine" re-usable, and display method flexible.
- No metrics. Use clean logs over metrics. Consider using metrics if this app becomes a long-running application.
- No config file. Avoid using config file util this app requires something meaningful configuration. Use cleanly defined hardcoded values for now.
- Input data is case-sensitive. This makes parsing simpler.
- No command is an error. This makes command handling simpler.
