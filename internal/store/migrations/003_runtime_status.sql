CREATE TABLE IF NOT EXISTS runtime_status (
  site_id TEXT PRIMARY KEY REFERENCES sites(id) ON DELETE CASCADE,
  org_id TEXT NOT NULL,
  deployment_id TEXT NOT NULL DEFAULT '',
  hosts TEXT[] NOT NULL DEFAULT '{}',
  status TEXT NOT NULL,
  started_at TIMESTAMPTZ,
  last_error TEXT NOT NULL DEFAULT '',
  requests_total BIGINT NOT NULL DEFAULT 0,
  errors_total BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_runtime_status_org ON runtime_status(org_id, status);
