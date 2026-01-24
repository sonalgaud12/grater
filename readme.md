```bash 
Note: This project is my propsal for open telemetry project for lfx spring'26 ( you can find the listed issue here: open-telemetry/opentelemetry-go-build-tools/issues/1528 )
```

### Grater

Analyze the impact of Go library changes on downstream dependencies by running their test suites.

### What is this?

Grater runs tests for packages that depend on your library. It compares results between versions so you can see what breaks before you release.

Tools like apidiff catch API changes. Grater catches behavior changes the stuff that breaks tests even when APIs stay compatible.

### Why?

Making changes to widely used Go libraries is stressful. You never know what will break downstream until someone files an issue.

This is especially true for foundational packages like OpenTelemetry that other CNCF projects

Grater gives you data instead of guesses.

### How it works

```bash
# Find packages that depend on yours
impact gather github.com/yourorg/yourlib -o dependents.json

# Run their tests against v1.2.0
impact test github.com/yourorg/yourlib v1.2.0 -o old.json

# Run their tests against v1.3.0
impact test github.com/yourorg/yourlib v1.3.0 -o new.json

# See what changed
impact compare old.json new.json
```

The comparison shows you:

Tests that started failing (your change probably broke something)

Tests that started passing (nice!)

Tests still failing (not your problem)

Tests still passing (you're good)

### Current state

This works but has some rough edges:

What works:

Gathers dependency lists

Runs test suites

Compares results between versions

Outputs JSON and readable summaries

What doesn't:

Tests run on your machine without isolation (security issue)

Tests run one at a time (slow)

The pkg.go.dev client is a stub (needs real API calls)

No caching (runs everything every time)

Can't resume if it crashes

### Installation

```bash
go install github.com/yourorg/grater/cmd/grater@latest
```

### Examples

Test a release candidate:

```bash
impact gather go.opentelemetry.io/otel --max 50
impact test go.opentelemetry.io/otel v1.21.0 -o stable.json
impact test go.opentelemetry.io/otel v1.22.0-rc1 -o rc.json
impact compare stable.json rc.json
```

Increase timeout for slow tests:

```bash
impact test mymodule v2.0.0 --timeout 600
```

### What needs work

**Security:**
Tests currently run directly on your machine. Anyone could put malicious code in their test suite. Need to run everything in containers with resource limits.

**Speed:**
One test at a time is painful. Should run many in parallel with a worker pool.

**API integration:**
Right now the dependency gatherer doesn't actually call pkg.go.dev. Need to implement the real API.

**Caching:**
Running the same tests repeatedly wastes time. Should cache results keyed by dependency version + target version.

**Better errors:**
Network hiccups kill the whole run. Need retries and the ability to resume.

### Planned improvements

**Next up:**

* Run tests in Docker containers
* Parallel test execution
* Real deps.dev API integration
* Result caching
* Retry logic

**After that:**

* Dependency graph analysis (only test what could be affected)
* GitHub Actions integration
* Config file support
* Smarter filtering 

**Eventually:**


* OpenTelemetry instrumentation
* Flaky test detection
* Impact predictions

### Project structure

```
grater/
├── cmd/impact/          # CLI
├── internal/
│   ├── compare/         # Compare results
│   ├── gather/          # Find dependencies
│   ├── models/          # Data types
│   └── test/            # Run tests
└── pkg/
    ├── client/          # API client
    └── report/          # Output formatting
```

### Who this helps

Library maintainers: Know what breaks before users complain

Teams with big codebases: Test dependency upgrades automatically

OSS projects: Show the community you care about compatibility


### Inspired by

Rust has Crater which does this for the whole ecosystem. This is the Go version, focused on individual libraries.

### Status

Early stage. Basic workflow works but not production ready. 
