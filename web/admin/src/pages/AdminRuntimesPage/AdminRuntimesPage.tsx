import { useState } from 'react';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { ConfirmActionDialog } from '../../components/molecules';
import { AdminRuntimeTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useGetAdminRuntimeSummaryQuery, useRestartAdminRuntimeMutation, useStopAdminRuntimeMutation } from '../../services/goGoHostApi';
import type { RuntimeStatus } from '../../services/types';
import './AdminRuntimesPage.css';

export function AdminRuntimesPage() {
  const summary = useGetAdminRuntimeSummaryQuery(undefined, { pollingInterval: 10_000 });
  const [restartRuntime] = useRestartAdminRuntimeMutation();
  const [stopRuntime] = useStopAdminRuntimeMutation();
  const [pending, setPending] = useState<{ action: 'restart' | 'stop'; runtime: RuntimeStatus } | undefined>();
  const [busySiteId, setBusySiteId] = useState<string | undefined>();
  async function confirmAction() {
    if (!pending) return;
    setBusySiteId(pending.runtime.siteId);
    try {
      if (pending.action === 'restart') await restartRuntime(pending.runtime.siteId).unwrap();
      else await stopRuntime(pending.runtime.siteId).unwrap();
      setPending(undefined);
    } finally {
      setBusySiteId(undefined);
    }
  }
  if (summary.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (summary.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load runtimes" error={apiErrorMessage(summary.error)} /></section>;
  const data = summary.data ?? { activeSites: 0, hosts: [], runtimes: [] };
  return (
    <section className="dashboard-panel admin-runtimes-page">
      <header>
        <div><h1>Runtimes</h1><p>Live and recently known runtime states across the host.</p></div>
        <button type="button" onClick={() => summary.refetch()}>Refresh</button>
      </header>
      <AdminRuntimeTable runtimes={data.runtimes} actionBusySiteId={busySiteId} onRestart={(runtime) => setPending({ action: 'restart', runtime })} onStop={(runtime) => setPending({ action: 'stop', runtime })} />
      <ConfirmActionDialog open={!!pending} title={pending?.action === 'restart' ? 'Restart runtime?' : 'Stop runtime?'} body={pending ? `${pending.action === 'restart' ? 'Restart' : 'Stop'} runtime for site ${pending.runtime.siteId}. This is a platform-wide operator action and will be audited.` : ''} confirmLabel={pending?.action === 'restart' ? 'Restart runtime' : 'Stop runtime'} busy={!!busySiteId} onCancel={() => setPending(undefined)} onConfirm={confirmAction} />
    </section>
  );
}
