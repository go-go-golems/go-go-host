import { EmptyState } from '../../components/atoms';
import './NoOrgsPage.css';

export function NoOrgsPage() {
  return <main className="no-orgs-page"><section className="dashboard-panel"><EmptyState title="Welcome to go-go-host" body="You do not belong to any organizations yet. Organization creation UI will be wired in the next increment." action={<button type="button" data-part="btn">Create organization</button>} /></section></main>;
}
