import { Outlet, useNavigate } from 'react-router-dom';
import { AppShell, AdminSidebar, type AdminSection } from '../../../components/organisms';
import { LoadingBlock } from '../../../components/atoms';
import { useGetConfigQuery, useGetMeQuery } from '../../../services/goGoHostApi';
import { isOIDCEnabled, logout } from '../../../auth/oidc';

function sectionFromPath(path: string): AdminSection {
  if (path.includes('/runtimes')) return 'runtimes';
  if (path.includes('/orgs')) return 'orgs';
  if (path.includes('/users')) return 'users';
  if (path.includes('/sites')) return 'sites';
  if (path.includes('/deployments')) return 'deployments';
  if (path.includes('/agents')) return 'agents';
  if (path.includes('/audit')) return 'audit';
  if (path.includes('/quotas')) return 'quotas';
  if (path.includes('/capabilities')) return 'capabilities';
  if (path.includes('/domains')) return 'domains';
  return 'overview';
}

export function AdminLayout() {
  const navigate = useNavigate();
  const me = useGetMeQuery();
  const config = useGetConfigQuery();
  if (me.isLoading) return <LoadingBlock lines={4} />;
  const active = sectionFromPath(location.pathname);
  const userLabel = `${me.data?.user.email ?? 'unknown user'} · platform admin`;
  return <AppShell memberships={me.data?.memberships ?? []} userLabel={userLabel} devAuth={config.data?.devAuth} onLogout={isOIDCEnabled(config.data) ? () => { void logout(config.data); } : undefined} sidebar={<AdminSidebar active={active} onSelect={(section) => navigate(`/admin/${section}`)} />}><Outlet /></AppShell>;
}
