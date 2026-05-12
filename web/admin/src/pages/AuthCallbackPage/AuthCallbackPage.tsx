import { useEffect, useState } from 'react';
import { LoadingBlock, ErrorCallout } from '../../components/atoms';
import { useGetConfigQuery } from '../../services/goGoHostApi';
import { completeLogin, isOIDCEnabled } from '../../auth/oidc';

export function AuthCallbackPage() {
  const config = useGetConfigQuery();
  const [error, setError] = useState('');

  useEffect(() => {
    if (!config.data) return;
    if (!isOIDCEnabled(config.data)) {
      setError('OIDC login is not enabled for this go-go-host instance.');
      return;
    }
    completeLogin(config.data)
      .then((returnTo) => { window.location.replace(returnTo); })
      .catch((err) => { setError(err instanceof Error ? err.message : String(err)); });
  }, [config.data]);

  if (config.isLoading || (!error && config.data)) return <LoadingBlock lines={4} />;
  if (config.error) return <ErrorCallout title="Unable to load auth config" error="The dashboard could not load /api/v1/config." />;
  return <ErrorCallout title="OIDC login failed" error={error || 'Unknown login callback error.'} onRetry={() => { window.location.assign('/app'); }} retryLabel="Return to dashboard" />;
}
