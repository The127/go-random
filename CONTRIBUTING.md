# Contributing to go-random

## Set up the tools

If you use Nix and direnv, run `direnv allow` in the repository. Without direnv, run `nix develop`. Either gives you a shell with Go, golangci-lint, govulncheck, just, lefthook and reuse, and installs the git hooks.

If you don't use Nix, install those tools yourself and run `lefthook install` once.

`just` lists the recipes. `just ci` runs every check that CI runs: the linters, the tests with the race detector, govulncheck and reuse. Run it before you open a pull request. The git hooks run the linters and the tests before each commit.

## Sign off your commits

By signing off a commit, you certify the [Developer Certificate of Origin](https://developercertificate.org/): that you wrote the change or have the right to submit it under the project's license. Add `-s` to `git commit`:

```sh
git commit -s
```

The commit-msg hook and CI refuse commits without a sign-off.

## Write commit messages and pull request titles

Use [Conventional Commits](https://www.conventionalcommits.org/), for example `fix: describe the change`. The types are `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `perf`, `build`, `ci` and `revert`.

Pull requests are squash merged, and the title becomes the commit message on `main`. CI checks that it follows the same rule.

If an AI tool helped write a commit, add an `Assisted-by` line with the tool and the exact model, like `Assisted-by: Claude:claude-opus-5-5`.

## Releases

[release-please](https://github.com/googleapis/release-please) reads the commit messages on `main` and keeps a release pull request open. When it merges, release-please tags the version and writes the changelog. A `fix` raises the patch version and a `feat` the minor version.
