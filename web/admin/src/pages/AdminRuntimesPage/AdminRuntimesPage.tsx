import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { AdminRuntimeTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useGetAdminRuntimeSummaryQuery } from '../../services/goGoHostApi';
import './AdminRuntimesPage.css';

export function AdminRuntimesPage() {
  const summary = useGetAdminRuntimeSummaryQuery(undefined, { pollingInterval: 10_000 });
  if (summary.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={6} /></section>;
  if (summary.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load runtimes" error={apiErrorMessage(summary.error)} /></section>;
  const data = summary.data ?? { activeSites: 0, hosts: [], runtimes: [] };
  return (
    <section className="dashboard-panel admin-runtimes-page">
      <header>
        <div><h1>Runtimes</h1><p>Live and recently known runtime states across the host.</p></div>
        <button type="button" onClick={() => summary.refetch()}>Refresh</button>
      </header>
      <AdminRuntimeTable runtimes={data.runtimes} />
    </section>
  );
}
