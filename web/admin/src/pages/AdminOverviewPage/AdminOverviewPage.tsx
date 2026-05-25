import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { MetricCard } from '../../components/molecules';
import { AdminRuntimeTable } from '../../components/organisms';
import { apiErrorMessage } from '../../services/errors';
import { useGetAdminRuntimeSummaryQuery } from '../../services/goGoHostApi';
import './AdminOverviewPage.css';

export function AdminOverviewPage() {
  const summary = useGetAdminRuntimeSummaryQuery(undefined, { pollingInterval: 10_000 });
  if (summary.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={5} /></section>;
  if (summary.error) return <section className="dashboard-panel"><ErrorCallout title="Unable to load platform runtime summary" error={apiErrorMessage(summary.error)} /></section>;
  const data = summary.data ?? { activeSites: 0, hosts: [], runtimes: [] };
  const requests = data.runtimes.reduce((total, runtime) => total + (runtime.requestsTotal ?? 0), 0);
  const errors = data.runtimes.reduce((total, runtime) => total + (runtime.errorsTotal ?? 0), 0);
  const failed = data.runtimes.filter((runtime) => runtime.status === 'failed').length;
  return (
    <div className="admin-overview-page">
      <section className="dashboard-panel admin-overview-page__hero">
        <h1>Platform admin</h1>
        <p>Global control-room view across go-go-host tenants, runtimes, deployments, and future policy surfaces.</p>
      </section>
      <section className="admin-overview-page__metrics" aria-label="Platform runtime summary">
        <MetricCard label="Active sites" value={data.activeSites} detail="Supervisor live runtimes" tone={data.activeSites ? 'success' : 'default'} />
        <MetricCard label="Known runtimes" value={data.runtimes.length} detail="Including stopped/failed states" />
        <MetricCard label="Hosts" value={data.hosts.length} detail={data.hosts.slice(0, 2).join(', ') || 'No hosts'} />
        <MetricCard label="Requests" value={requests} detail={`${errors} errors`} tone={errors ? 'danger' : 'default'} />
        <MetricCard label="Failed" value={failed} detail="Runtime failures" tone={failed ? 'danger' : 'success'} />
      </section>
      <section className="dashboard-panel">
        <header className="admin-overview-page__section-header"><h2>Runtime snapshot</h2><a href="/admin/runtimes">Open runtimes</a></header>
        <AdminRuntimeTable runtimes={data.runtimes.slice(0, 5)} />
      </section>
    </div>
  );
}
