ALTER TABLE deployments
  ADD COLUMN IF NOT EXISTS bundle_sha256 TEXT NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS runtime_events (
  id TEXT PRIMARY KEY,
  site_id TEXT NOT NULL DEFAULT '',
  org_id TEXT NOT NULL DEFAULT '',
  deployment_id TEXT NOT NULL DEFAULT '',
  event_type TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT '',
  message TEXT NOT NULL DEFAULT '',
  metadata_json JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_runtime_events_site_created ON runtime_events(site_id, created_at DESC);
