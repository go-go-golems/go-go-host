---
Title: Hosted Capability Security Review Checklist
Ticket: HOST-001-GO-GO-HOST-V1
DocType: reference
Topics: [go-go-host, security, hardening]
Status: active
Intent: operational
---

# Hosted Capability Security Review Checklist

Use this checklist before adding or enabling any new capability exposed to hosted Goja code.

## Required questions

- What exact JavaScript API is exposed?
- Is the API deterministic and tenant-scoped?
- Can the API read host files, process environment, network sockets, or subprocesses?
- Does the API need quota, timeout, or byte-limit enforcement?
- Is access controlled by `site_capabilities`?
- Does deployment validation report requested vs effective capability state?
- Are denied capabilities tested?
- Are runtime failures observable through runtime events/status/audit?

## Default deny

The following remain unavailable in hosted v1:

- unrestricted `fs`,
- `exec`/subprocess execution,
- process environment passthrough,
- plaintext secret reads,
- arbitrary outbound network access.

## Test requirements

Every new capability needs:

- deployment validation test for denied policy,
- runtime integration test for allowed behavior,
- quota or timeout test if the capability consumes resources,
- dashboard/admin visibility if operators can toggle it.
