import type { ReactNode } from 'react';
import { Navigate, useParams } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { useGetConfigQuery, useGetMeQuery } from '../../services/goGoHostApi';
import { beginLogin, clearTokens, isOIDCEnabled } from '../../auth/oidc';

export function RequireSession({ children }: { children: ReactNode }) {
  const config = useGetConfigQuery();
  const me = useGetMeQuery();
  if (config.isLoading || me.isLoading) return <LoadingBlock lines={4} />;
  if ((me.error || !me.data) && isOIDCEnabled(config.data)) {
    void beginLogin(config.data, window.location.pathname + window.location.search);
    return <LoadingBlock lines={4} />;
  }
  if (me.error || !me.data) return <ErrorCallout title="Unable to load session" error="The dashboard could not load /api/v1/me." />;
  return <>{children}</>;
}

export function RequirePlatformAdmin({ children }: { children: ReactNode }) {
  const me = useGetMeQuery();
  if (me.isLoading) return <LoadingBlock lines={4} />;
  if (me.error || !me.data) return <ErrorCallout title="Unable to load platform session" error="The admin dashboard could not load /api/v1/me." onRetry={() => { clearTokens(); window.location.assign('/app'); }} retryLabel="Sign in again" />;
  if (!me.data.platformAdmin) return <ErrorCallout title="Platform admin required" error="Current user is authenticated but is not allowed to inspect platform-wide state." />;
  return <>{children}</>;
}

export function RequireOrgAccess({ children }: { children: ReactNode }) {
  const { orgId } = useParams();
  const me = useGetMeQuery();
  if (me.isLoading) return <LoadingBlock lines={4} />;
  if (me.error || !me.data) return <ErrorCallout title="Unable to load organizations" error="The dashboard could not load organization memberships." />;
  if (!orgId) return <Navigate to="/app" replace />;
  const membership = me.data.memberships.find((m) => m.orgId === orgId);
  if (!membership) return <ErrorCallout title="Organization access denied" error={`Current user is not a member of ${orgId}.`} />;
  return <>{children}</>;
}
