# Docmgr and ticket workflow

Use docmgr tickets for work that requires investigation, design discussion, implementation history, screenshots, or a durable diary. Do not use tickets as a replacement for stable repository documentation. Stable contributor guidance belongs under `docs`.

## When to create a ticket

Create or reuse a docmgr ticket when the work involves:

- Architecture or design research.
- Multi-step implementation across subsystems.
- Debugging with non-obvious failures.
- UI screenshot review.
- Operational runbooks.
- Security, deployment, runtime, or production-readiness analysis.
- Work that another contributor may need to resume later.

Small code-only fixes do not always need a ticket. If the reasoning would be hard to reconstruct from code and tests, create one.

## Standard ticket setup

```bash
docmgr ticket create-ticket \
  --ticket HOST-123 \
  --title "Short descriptive title" \
  --topics go-go-host,documentation

docmgr doc add --ticket HOST-123 --doc-type design-doc --title "Design title"
docmgr doc add --ticket HOST-123 --doc-type reference --title "Investigation diary"
```

Typical ticket files:

```text
ttmp/YYYY/MM/DD/HOST-123--slug/
  index.md
  tasks.md
  changelog.md
  design-doc/
  reference/
  playbook/
  sources/
  scripts/
```

## Diary expectations

Keep a diary for non-trivial work. The diary should record:

- The goal.
- The prompt or task context.
- Commands run.
- Files inspected or changed.
- What worked.
- What failed, with exact errors.
- What was learned.
- What should be reviewed.
- Validation results.
- Follow-up work.

Update the diary while working, not only at the end. A useful diary lets another contributor resume the task without reading chat logs.

## Relating files

Relate important files to the relevant doc. Use absolute paths for clarity.

```bash
docmgr doc relate \
  --doc ttmp/YYYY/MM/DD/HOST-123--slug/design-doc/01-design.md \
  --file-note "/abs/path/to/file.go:Why this file matters"
```

Relate only meaningful files. Do not attach every file touched by a large refactor if only a few explain the design.

## Tasks and changelog

Use tasks for planned work:

```bash
docmgr task add --ticket HOST-123 --text "Implement backend API"
docmgr task check --ticket HOST-123 --id 1
```

Use changelog entries for completed steps:

```bash
docmgr changelog update --ticket HOST-123 \
  --entry "Implemented backend API and integration tests." \
  --file-note "/abs/path/to/file.go:Primary implementation"
```

## Stable docs vs ticket docs

Use this rule:

| Content | Location |
|---|---|
| Temporary investigation notes | `ttmp/.../reference` |
| Design alternatives and rationale | `ttmp/.../design-doc` |
| Ticket-specific run commands | `ttmp/.../playbook` |
| Reusable contributor workflow | `docs/contributing` |
| Reusable architecture reference | `docs/architecture` |
| Reusable local/operational commands | `docs/runbooks` |

Promote a ticket lesson into `docs` when it becomes a rule or workflow future contributors should follow.

## Validation

Run doctor before handoff:

```bash
docmgr doctor --ticket HOST-123 --stale-after 30
```

Fix vocabulary warnings by adding appropriate vocabulary entries, or change the topic/doc type to an existing value if the new term is unnecessary.

## Commit guidance

Commit at logical intervals:

1. Ticket setup and initial design package.
2. Stable docs or code implementation phase.
3. Validation and diary/changelog updates.

Keep unrelated worktree changes out of the commit. Always inspect `git status --short` before staging.
