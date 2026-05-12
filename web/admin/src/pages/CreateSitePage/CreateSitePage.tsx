import { useNavigate, useParams } from 'react-router-dom';
import './CreateSitePage.css';
export function CreateSitePage() {
  const { orgId } = useParams();
  const navigate = useNavigate();
  return <section className="dashboard-panel create-site-page"><h1>Create site</h1><p>This form shell is ready; mutation wiring comes in the next slice.</p><label>Slug <input placeholder="hello" /></label><label>Name <input placeholder="Hello Site" /></label><div className="create-site-page__actions"><button type="button" data-part="btn" onClick={() => navigate(`/app/orgs/${orgId}/sites`)}>Cancel</button><button type="button" data-part="btn">Create site</button></div></section>;
}
