# AGENTS.md

This file is for AI coding tools that work on go-random. `CLAUDE.md` is a symlink to it.

go-random is a small library that gives Go code cryptographically secure randomness through an interface. Keep it small and free of dependencies outside the standard library. Production code must never be able to get predictable randomness. That's why the mock stays in `randomtest` and panics outside a test binary.

Follow [CONTRIBUTING.md](CONTRIBUTING.md). Run `just ci` before you call a change done. If a change adds or changes behaviour, add tests for it. Give every commit you helped write an `Assisted-by` trailer with your name and exact model version, and don't add a `Co-Authored-By` line.
