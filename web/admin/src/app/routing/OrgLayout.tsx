import { Outlet, useNavigate, useParams } from 'react-router-dom';
import { AppShell, OrgSidebar, type OrgSection } from '../../components/organisms';
import { LoadingBlock } from '../../components/atoms';
import { useGetConfigQuery, useGetMeQuery } from '../../services/goGoHostApi';

export function OrgLayout() {
  const { orgId } = useParams();
  const navigate = useNavigate();
  const me = useGetMeQuery();
  const config = useGetConfigQuery();
  const path = location.pathname;
  const active: OrgSection = path.includes('/agents') ? 'agents' : path.includes('/audit') ? 'audit' : 'sites';
  if (me.isLoading) return <LoadingBlock lines={4} />;
  const userLabel = me.data?.user.email ?? 'unknown user';
  return <AppShell memberships={me.data?.memberships ?? []} selectedOrgId={orgId} userLabel={userLabel} devAuth={config.data?.devAuth} onOrgSelect={(id) => navigate(`/app/orgs/${id}/sites`)} sidebar={<OrgSidebar active={active} onSelect={(section) => navigate(`/app/orgs/${orgId}/${section}`)} />}><Outlet /></AppShell>;
}
