CREATE INDEX IF NOT EXISTS idx_memberships_user ON memberships(user_id, org_id);
CREATE INDEX IF NOT EXISTS idx_sites_org ON sites(org_id, slug);
CREATE INDEX IF NOT EXISTS idx_deployments_site_version ON deployments(site_id, version);
CREATE INDEX IF NOT EXISTS idx_deploy_runs_site ON deploy_runs(site_id, status);
CREATE INDEX IF NOT EXISTS idx_agents_org ON agents(org_id, status);
CREATE INDEX IF NOT EXISTS idx_audit_log_org_time ON audit_log(org_id, created_at);
CREATE INDEX IF NOT EXISTS idx_audit_log_resource ON audit_log(resource_type, resource_id, created_at);
