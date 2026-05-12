import { Navigate } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { useGetMeQuery } from '../../services/goGoHostApi';
import { NoOrgsPage } from '../../pages/NoOrgsPage';

export function OrgRedirectOrOnboarding() {
  const me = useGetMeQuery();
  if (me.isLoading) return <LoadingBlock lines={4} />;
  if (me.error || !me.data) return <ErrorCallout title="Unable to load session" error="The dashboard could not load /api/v1/me." />;
  if (me.data.memberships.length === 0) return <NoOrgsPage />;
  return <Navigate to={`/app/orgs/${me.data.memberships[0].orgId}/sites`} replace />;
}
