CREATE TABLE IF NOT EXISTS site_config (
  site_id TEXT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  key TEXT NOT NULL,
  value_json JSONB NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (site_id, key)
);

CREATE INDEX IF NOT EXISTS idx_site_config_site ON site_config(site_id);
CREATE INDEX IF NOT EXISTS idx_site_domains_site_status ON site_domains(site_id, status);
