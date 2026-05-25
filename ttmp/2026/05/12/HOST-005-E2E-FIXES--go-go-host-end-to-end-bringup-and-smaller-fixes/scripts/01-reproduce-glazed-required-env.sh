#!/usr/bin/env bash
set -euo pipefail

# Reproduce/confirm a Glazed parser issue: AppName-based env loading works for
# optional fields, but a required flag can fail during Cobra parsing before the
# env middleware has a chance to provide the value.
#
# Expected result against the current local Glazed checkout:
#   - optional env field parses successfully
#   - required env field returns "Field required-name is required"
#   - script exits 0 because the bug was reproduced

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORKSPACE_ROOT="$(cd "$SCRIPT_DIR/../../../../../../.." && pwd)"
GLAZED_DIR="$WORKSPACE_ROOT/glazed"

if [[ ! -d "$GLAZED_DIR" ]]; then
  echo "could not find local glazed checkout at $GLAZED_DIR" >&2
  exit 2
fi

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

cat >"$TMPDIR/go.mod" <<EOF_GO_MOD
module required-env-repro

go 1.26.1

require github.com/go-go-golems/glazed v1.2.7

replace github.com/go-go-golems/glazed => $GLAZED_DIR
EOF_GO_MOD

cat >"$TMPDIR/main.go" <<'EOF_GO'
package main

import (
	"fmt"
	"os"

	glazedcli "github.com/go-go-golems/glazed/pkg/cli"
	glazedcmds "github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
)

func parseField(name string, required bool) error {
	desc := glazedcmds.NewCommandDescription(
		"probe",
		glazedcmds.WithFlags(fields.New(
			name,
			fields.TypeString,
			fields.WithRequired(required),
			fields.WithHelp("field populated from env in this repro"),
		)),
	)
	parser, err := glazedcli.NewCobraParserFromSections(desc.Schema, &glazedcli.CobraParserConfig{
		ShortHelpSections: []string{schema.DefaultSlug},
		AppName:           "REQ_ENV_TEST",
	})
	if err != nil {
		return fmt.Errorf("new parser: %w", err)
	}
	cobraCmd := glazedcli.NewCobraCommandFromCommandDescription(desc)
	if err := parser.AddToCobraCommand(cobraCmd); err != nil {
		return fmt.Errorf("add to cobra: %w", err)
	}
	_, err = parser.Parse(cobraCmd, nil)
	return err
}

func main() {
	if err := os.Setenv("REQ_ENV_TEST_OPTIONAL_NAME", "from-env"); err != nil {
		panic(err)
	}
	if err := os.Setenv("REQ_ENV_TEST_REQUIRED_NAME", "from-env"); err != nil {
		panic(err)
	}

	if err := parseField("optional-name", false); err != nil {
		fmt.Printf("UNEXPECTED: optional env-backed field failed: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("OK: optional env-backed field parses")

	err := parseField("required-name", true)
	if err == nil {
		fmt.Println("BUG NOT REPRODUCED: required env-backed field parsed successfully")
		os.Exit(1)
	}

	fmt.Printf("BUG REPRODUCED: required env-backed field failed before env could satisfy it: %v\n", err)
}
EOF_GO

(
  cd "$TMPDIR"
  GOWORK=off go mod tidy >/dev/null
  GOWORK=off go run .
)
