import { Outlet, useNavigate, useParams } from 'react-router-dom';
import { AppShell, OrgSidebar, type OrgSection } from '../../components/organisms';
import { LoadingBlock } from '../../components/atoms';
import { useGetConfigQuery, useGetMeQuery } from '../../services/goGoHostApi';
import { isOIDCEnabled, logout } from '../../auth/oidc';

export function OrgLayout() {
  const { orgId } = useParams();
  const navigate = useNavigate();
  const me = useGetMeQuery();
  const config = useGetConfigQuery();
  const path = location.pathname;
  const active: OrgSection = path.includes('/docs') ? 'docs' : path.includes('/agents') ? 'agents' : path.includes('/audit') ? 'audit' : path.includes('/members') ? 'members' : path.includes('/usage') ? 'usage' : 'sites';
  if (me.isLoading) return <LoadingBlock lines={4} />;
  const userLabel = me.data?.user.email ?? 'unknown user';
  return <AppShell memberships={me.data?.memberships ?? []} selectedOrgId={orgId} userLabel={userLabel} devAuth={config.data?.devAuth} onLogout={isOIDCEnabled(config.data) ? () => { void logout(config.data); } : undefined} onOrgSelect={(id) => navigate(`/app/orgs/${id}/sites`)} sidebar={<OrgSidebar active={active} onSelect={(section) => navigate(`/app/orgs/${orgId}/${section}`)} />}><Outlet /></AppShell>;
}
