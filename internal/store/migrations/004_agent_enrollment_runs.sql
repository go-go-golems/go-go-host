CREATE TABLE IF NOT EXISTS agent_enrollment_tokens (
  token_hash TEXT PRIMARY KEY,
  agent_id TEXT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
  org_id TEXT NOT NULL REFERENCES orgs(id) ON DELETE CASCADE,
  status TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_agent_enrollment_tokens_agent ON agent_enrollment_tokens(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_keys_agent_status ON agent_keys(agent_id, status);
CREATE INDEX IF NOT EXISTS idx_deploy_runs_agent ON deploy_runs(agent_id, created_at DESC);
