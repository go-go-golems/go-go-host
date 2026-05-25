ALTER TABLE agent_site_grants
  ADD COLUMN IF NOT EXISTS can_activate BOOLEAN NOT NULL DEFAULT false;
