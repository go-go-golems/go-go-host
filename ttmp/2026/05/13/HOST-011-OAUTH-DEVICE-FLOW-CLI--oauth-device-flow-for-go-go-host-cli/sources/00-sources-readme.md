---
Title: HOST-011 Source Index
Ticket: HOST-011-OAUTH-DEVICE-FLOW-CLI
Status: active
Topics:
    - go-go-host
    - cli
    - oauth
    - oidc
    - keycloak
    - security
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Source material captured for HOST-011 OAuth Device Flow CLI design.
WhatFor: Preserve source evidence for the HOST-011 design guide.
WhenToUse: Use when reviewing or implementing OAuth Device Flow CLI support.
---

# HOST-011 Sources

This folder stores the primary external references and live observations used for the OAuth Device Flow CLI design. Markdown web-page captures were produced with `defuddle parse ... --md` unless noted otherwise.

## Protocol specifications

- `01-rfc8628-oauth-device-authorization-grant.md`
  - Source: https://datatracker.ietf.org/doc/html/rfc8628
  - Purpose: Primary specification for OAuth 2.0 Device Authorization Grant. Defines `device_code`, `user_code`, `verification_uri`, `verification_uri_complete`, polling, `authorization_pending`, `slow_down`, `access_denied`, and `expired_token`.

- `02-rfc8414-oauth-authorization-server-metadata.md`
  - Source: https://datatracker.ietf.org/doc/html/rfc8414
  - Purpose: Defines authorization-server metadata discovery. This explains why the CLI should discover `device_authorization_endpoint` and `token_endpoint` from `/.well-known/openid-configuration` instead of hard-coding Keycloak paths.

- `03-rfc7009-oauth-token-revocation.md`
  - Source: https://datatracker.ietf.org/doc/html/rfc7009
  - Purpose: Reference for a future `go-go-host logout` implementation that revokes refresh tokens instead of only deleting local config.

- `04-rfc8252-oauth-native-apps.md`
  - Source: https://datatracker.ietf.org/doc/html/rfc8252
  - Purpose: Reference for the alternative native-app authorization-code flow with loopback redirect. Useful for explaining why device flow is simpler for CLI login and what tradeoffs it has.

## Keycloak references

- `05-keycloak-oidc-layers.md`
  - Source: https://www.keycloak.org/securing-apps/oidc-layers
  - Purpose: Official Keycloak endpoint reference. Confirms the device authorization endpoint path `/realms/{realm-name}/protocol/openid-connect/auth/device`, that public clients can invoke it, and that Device Authorization Grant is preferred over password/direct grant for this use case.

- `06-keycloak-community-device-grant-design.md`
  - Source: https://github.com/keycloak/keycloak-community/blob/main/design/oauth2-device-authorization-grant.md
  - Purpose: Keycloak design note for device grant support. Documents Keycloak-specific behavior, the short verification endpoint, default device-code lifespan/polling interval concepts, and the per-client setting named "OAuth 2.0 Device Authorization Grant Enabled".

- `07-terraform-provider-keycloak-openid-client-defuddle.md`
  - Source: https://github.com/keycloak/terraform-provider-keycloak/blob/main/docs/resources/openid_client.md
  - Purpose: Terraform provider documentation captured through Defuddle. Used to confirm that `keycloak_openid_client` exposes device-flow settings.

- `08-terraform-provider-keycloak-openid-client.md`
  - Source: https://raw.githubusercontent.com/keycloak/terraform-provider-keycloak/master/docs/resources/openid_client.md
  - Purpose: Raw Markdown fallback for the Terraform provider docs. This file is easier to grep than the GitHub-rendered Defuddle capture. It confirms `oauth2_device_authorization_grant_enabled`, `oauth2_device_code_lifespan`, and `oauth2_device_polling_interval`.

## Live production observations

- `live-go-go-host-openid-configuration.json`
  - Source: https://auth.yolo.scapegoat.dev/realms/go-go-host/.well-known/openid-configuration
  - Purpose: Live OIDC discovery document for the production realm. Confirms that the realm advertises `device_authorization_endpoint` and supports the `urn:ietf:params:oauth:grant-type:device_code` grant type at the metadata level.

- `live-device-endpoint-disabled-response.json`
  - Source: POST to `https://auth.yolo.scapegoat.dev/realms/go-go-host/protocol/openid-connect/auth/device` with `client_id=go-go-host-dashboard`.
  - Purpose: Live negative evidence. Production Keycloak has the realm-level endpoint, but the current dashboard client is not yet allowed to initiate Device Authorization Grant. The response is `unauthorized_client` with the message that the flow is disabled for the client.
