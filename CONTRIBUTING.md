# Contributing to Heimdall

Heimdall is the shared configuration library for the Asgard EHS ecosystem.
Odin, Muninn, Huginn, and Ratatoskr import it directly; the CLI is a thin
wrapper over the same library. Because every tool in the family depends on
it, Heimdall is held to a higher stability bar than most libraries its size.

Contributions that keep the library small, validated, and predictable are
welcome. Contributions that expand scope, loosen validation, or introduce
new coupling between tools will be discussed before they are written, not
after. This document is direct about that because the alternative is wasted
effort. Your time is worth respecting, and so is the project's.

## Before You Start

A few things to know about how this project is run, so you can decide
whether contributing here is a good use of your time.

Heimdall is solo-maintained. Review cadence depends on the maintainer's
availability, and not every contribution will be accepted. The project has
a defined scope and a defined shape, and changes that don't fit either will
be declined — sometimes after discussion, sometimes on sight. Declining a
contribution is not a comment on its quality; it means it doesn't belong in
this particular library.

Heimdall sits at the center of the Asgard EHS dependency graph. A change
that makes Heimdall more flexible for one caller often makes it less
predictable for the others. The maintainer's bias is toward keeping the
library boring: small API surface, strict validation, no surprises. Changes
that support that orientation are welcome. Changes that trade predictability
for cleverness usually are not.

## Reporting Issues

A good bug report for Heimdall includes, at minimum:

- Your operating system (platform-specific database paths are a common source
  of issues)
- Your Go version (`go version`)
- The Heimdall version or commit you are using
- The namespace and key involved, if the bug is about a specific config entry
- A minimal reproduction — ideally a short Go file or a sequence of CLI
  commands that demonstrates the problem, along with the output you saw

If the bug involves the SQLite database itself, please mention whether the
database was freshly seeded or migrated from an earlier Heimdall version,
since that distinction changes how the bug should be investigated.

If you are not sure whether something is a bug or expected behavior, open an
issue anyway and label it as a question. An unclear behavior that needs
documentation is itself a kind of bug.

## Before You Submit a PR

There are two kinds of pull requests: ones that should just be sent, and
ones that should be discussed first. Knowing which is which saves
everyone's time.

**Just send it.** Typo fixes, broken-link corrections, obvious-bug fixes
with a clear root cause, small documentation improvements, and test
additions for existing behavior don't need prior discussion. Open a PR.

**Open an issue first.** Anything that fits one of the following categories
should be discussed before code is written:

- Changes to the Go API (`Open`, `Get`, `Set`, `List`, `Reset`, `Schema`,
  `ChangeNotifier`, or any other exported symbol)
- Changes to config schemas — adding a namespace, adding a key, changing a
  default, changing a type, or changing validation rules (see below for the
  specific process)
- Changes to the CLI surface — new subcommands, renamed flags, changes to
  output format that scripts might parse
- New database migrations or changes to existing ones
- New validation types beyond the current set (`string`, `path`, `secret`,
  `number`, `boolean`, `array`, `enum`)
- Refactors that touch more than a handful of files
- Anything you're not sure about

Pull requests in these categories that arrive without prior discussion will
usually be closed without a detailed review. This is not personal, and it
is not a comment on code quality. It is a consequence of scarce review time
being spent on changes whose direction has already been agreed on. A
ten-minute issue conversation before you start can save you hours of work
on a PR that does not fit the project's direction.

### Config Schema Changes

Config schemas are Heimdall's real public API — more than the Go types, they
are what downstream tools actually depend on. A misplaced schema change can
silently break every consumer of the library. Because of that, schema
changes have a specific process.

**New config keys originate with the consuming tool, not with Heimdall.**
If Muninn needs a new setting, the discussion about what the setting should
be named, typed, and defaulted to happens on the Muninn repository. Once
the consuming tool has decided what it needs, a mechanical PR against
Heimdall adds the schema entry. This keeps the design conversation where
the context is — the team that understands why the setting is needed — and
keeps Heimdall's review scope to "does this schema entry look correct."

For changes to existing schemas (altering a default, tightening validation,
renaming a key), open the issue on the Heimdall repository directly, since
every consuming tool is affected and the discussion belongs in one place.

Schema changes should always include a migration path for existing
databases. A new key with a default is safe; renaming an existing key
requires the old key to be readable during a transition period, not
deleted outright.

### Database Migrations

Heimdall's SQLite database has a migration history that must be
forward-only in practice. Once a migration has shipped in a released
version, it cannot be edited — it can only be superseded by a later
migration. New migrations must be safe against existing data, and they
must be tested against a database that was created by an earlier version
of Heimdall, not just against a fresh database.

Migration numbering follows the existing convention. Do not reuse or
rearrange numbers.

## Development Setup

```bash
git clone https://github.com/asgardehs/heimdall.git
cd heimdall
go build ./cmd/heimdall
go test ./...
```

To run the CLI against a test database without touching your real
configuration, set `HEIMDALL_DATA_DIR` to a scratch directory:

```bash
HEIMDALL_DATA_DIR=/tmp/heimdall-test ./heimdall config list odin
```

The scratch directory is created on first use and can be deleted at any
time. Regenerating it exercises the initial migration and seed path, which
is useful when working on those code paths.

## What Makes a Good PR

A pull request is more likely to be accepted quickly if it:

- Has a clear, narrow scope — one change per PR, not five
- Preserves the no-CGO dependency posture (Heimdall uses
  `modernc.org/sqlite` for a reason)
- Keeps the Go API backwards-compatible, unless the change is part of an
  agreed major version bump
- Keeps the CLI output format stable, unless the change is explicitly
  approved as a breaking change to the CLI
- Includes tests for behavior changes, not just code changes
- Includes tests against the CLI for CLI changes, not just the library API
- Updates documentation when it changes public surface, including the
  hosted docs site pages when those are affected

A pull request that bundles a bug fix with a refactor with a new feature
will be asked to split before it is reviewed. Not because the work is
unwelcome, but because reviewing bundled changes is harder and slower, and
each piece deserves its own decision.

## Commits and PRs

Commit messages follow the Asgard EHS project conventions, which are also
documented in the
[brand guidelines](https://asgardehs.github.io/docs/brand/#voice-and-tone):

- Imperative mood: "Add array type validation," not "Added" or "Adds"
- Present tense
- Summary line under 72 characters
- If the change is not self-explanatory, a blank line followed by a
  paragraph that explains why — not what, since the diff shows what

Pull request descriptions should link to the issue they address (if one
exists), summarize the change in one or two sentences, and call out
anything a reviewer should pay particular attention to. PRs that are
simply labeled "fix bug" or "update code" will be asked for more context
before review.

## Code of Conduct

Participation in this project is governed by the
[Asgard EHS Code of Conduct](CODE_OF_CONDUCT.md). By contributing, you agree
to uphold its expectations.

## License

Heimdall is licensed under Apache-2.0. Contributions are accepted under the
same license. By submitting a pull request, you confirm that you have the
right to submit the code and that you agree to license it under Apache-2.0.
