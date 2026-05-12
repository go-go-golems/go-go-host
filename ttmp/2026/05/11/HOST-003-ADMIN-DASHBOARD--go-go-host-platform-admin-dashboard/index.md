---
Title: go-go-host platform admin dashboard
Ticket: HOST-003-ADMIN-DASHBOARD
Status: active
Topics:
    - dashboard
    - frontend
    - go-go-host
    - rtk-query
    - storybook
    - platform-admin
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: internal/httpapi/admin_inventory.go
      Note: Admin org/user/site/deployment HTTP inventory handlers
    - Path: internal/httpapi/runtime.go
      Note: Backend platform-admin-gated runtime summary endpoint
    - Path: internal/store/admin.go
      Note: Store wrappers for admin inventory rows
    - Path: web/admin/src/app/routes.tsx
      Note: Defines /admin route tree alongside /app
    - Path: web/admin/src/pages/AdminDeploymentsPage/AdminDeploymentsPage.tsx
      Note: Admin deployment inventory page
    - Path: web/admin/src/services/goGoHostApi.ts
      Note: RTK Query endpoint for admin runtime summary
ExternalSources: []
Summary: ""
LastUpdated: 2026-05-11T22:14:39.427159207-04:00
WhatFor: ""
WhenToUse: ""
---



# go-go-host platform admin dashboard

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- dashboard
- frontend
- go-go-host
- rtk-query
- storybook
- platform-admin

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
