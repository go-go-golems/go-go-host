import { type FormEvent, useId, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ErrorCallout, LoadingBlock } from '../../components/atoms';
import { SiteHostCopy } from '../../components/molecules';
import { apiErrorMessage } from '../../services/errors';
import { useCreateSiteMutation, useGetConfigQuery } from '../../services/goGoHostApi';
import './CreateSitePage.css';

const slugPattern = /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/;

function validateSite(slug: string, name: string): string[] {
  const errors: string[] = [];
  if (!slug.trim()) errors.push('Slug is required.');
  if (slug && !slugPattern.test(slug)) errors.push('Slug must be lowercase DNS-safe text: letters, numbers, and single hyphens, without leading or trailing hyphens.');
  if (!name.trim()) errors.push('Name is required.');
  if (name.length > 120) errors.push('Name must be 120 characters or fewer.');
  return errors;
}

export function CreateSitePage() {
  const { orgId } = useParams();
  const navigate = useNavigate();
  const slugId = useId();
  const nameId = useId();
  const config = useGetConfigQuery();
  const [createSite, createState] = useCreateSiteMutation();
  const [slug, setSlug] = useState('');
  const [name, setName] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const normalizedSlug = slug.trim().toLowerCase();
  const normalizedName = name.trim();
  const errors = useMemo(() => validateSite(normalizedSlug, normalizedName), [normalizedSlug, normalizedName]);
  const previewHost = normalizedSlug && config.data?.baseDomain ? `${normalizedSlug}.${config.data.baseDomain}` : '';

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitted(true);
    if (!orgId || errors.length > 0) return;
    const site = await createSite({ orgId, slug: normalizedSlug, name: normalizedName }).unwrap();
    navigate(`/app/orgs/${orgId}/sites/${site.id}`);
  }

  if (config.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={4} /></section>;

  return <section className="dashboard-panel create-site-page">
    <header>
      <h1>Create site</h1>
      <p>Pick a DNS-safe slug. You can upload and activate a deployment after the site exists.</p>
    </header>
    {submitted && errors.length > 0 ? <ErrorCallout title="Fix site details" error={errors.join('\n')} /> : null}
    {createState.error ? <ErrorCallout title="Unable to create site" error={apiErrorMessage(createState.error)} /> : null}
    <form onSubmit={onSubmit} noValidate>
      <label htmlFor={slugId}>Slug</label>
      <input id={slugId} value={slug} onChange={(event) => setSlug(event.target.value)} placeholder="hello-world" aria-invalid={submitted && errors.some((e) => e.startsWith('Slug'))} />
      <small>Lowercase letters, numbers, and hyphens. This becomes the default host prefix.</small>

      <label htmlFor={nameId}>Name</label>
      <input id={nameId} value={name} onChange={(event) => setName(event.target.value)} placeholder="Hello World" aria-invalid={submitted && errors.some((e) => e.startsWith('Name'))} />
      <small>A human-readable name for operators and developers.</small>

      {previewHost ? <div className="create-site-page__preview"><span>Preview host</span><SiteHostCopy host={previewHost} /></div> : null}

      <div className="create-site-page__actions">
        <button type="button" data-part="btn" onClick={() => navigate(`/app/orgs/${orgId}/sites`)}>Cancel</button>
        <button type="submit" data-part="btn" disabled={createState.isLoading}>{createState.isLoading ? 'Creating…' : 'Create site'}</button>
      </div>
    </form>
  </section>;
}
