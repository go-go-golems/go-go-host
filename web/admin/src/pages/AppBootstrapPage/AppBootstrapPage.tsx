import { StatusPill } from '../../components/atoms/StatusPill';
import { useGetConfigQuery, useGetMeQuery, useListSitesQuery } from '../../services/goGoHostApi';
import './AppBootstrapPage.css';

export function AppBootstrapPage() {
  const config = useGetConfigQuery();
  const me = useGetMeQuery();
  const selectedOrg = me.data?.memberships[0];
  const sites = useListSitesQuery(selectedOrg?.orgId ?? '', { skip: !selectedOrg });

  if (config.isLoading || me.isLoading) {
    return <main className="dashboard-shell"><section className="dashboard-panel">Loading session…</section></main>;
  }
  if (config.error || me.error) {
    return <main className="dashboard-shell"><section className="dashboard-panel dashboard-error">Unable to load dashboard session.</section></main>;
  }
  if (!selectedOrg) {
    return <main className="dashboard-shell"><section className="dashboard-panel"><h1>Welcome to go-go-host</h1><p>You do not belong to any organizations yet.</p></section></main>;
  }

  return (
    <main className="dashboard-shell">
      <header className="dashboard-topbar">
        <strong>go-go-host</strong>
        <span>Org: {selectedOrg.orgName}</span>
        {config.data?.devAuth ? <StatusPill status="dev auth" tone="info" /> : null}
      </header>
      <section className="dashboard-panel">
        <h1>Sites</h1>
        <p className="muted">Signed in as {me.data?.user.email}. This is the first dashboard scaffold backed by RTK Query and MSW stories.</p>
        {sites.isLoading ? <p>Loading sites…</p> : null}
        {sites.error ? <p className="dashboard-error">Unable to load sites.</p> : null}
        <div className="site-grid">
          {sites.data?.map((site) => (
            <article className="site-card" key={site.id}>
              <h2>{site.name}</h2>
              <p>{site.primaryHost}</p>
              <StatusPill status={site.status} />
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}
