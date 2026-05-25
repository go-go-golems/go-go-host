import { type FormEvent, useId, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { EmptyState, ErrorCallout } from '../../components/atoms';
import { apiErrorMessage } from '../../services/errors';
import { useCreateOrgMutation } from '../../services/goGoHostApi';
import './NoOrgsPage.css';

const slugPattern = /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/;
function validateOrg(slug: string, name: string): string[] {
  const errors: string[] = [];
  if (!slug) errors.push('Organization slug is required.');
  if (slug && !slugPattern.test(slug)) errors.push('Organization slug must be lowercase DNS-safe text.');
  if (!name) errors.push('Organization name is required.');
  return errors;
}

export function NoOrgsPage() {
  const navigate = useNavigate();
  const slugId = useId();
  const nameId = useId();
  const [slug, setSlug] = useState('');
  const [name, setName] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const [createOrg, createState] = useCreateOrgMutation();
  const normalizedSlug = slug.trim().toLowerCase();
  const normalizedName = name.trim();
  const errors = useMemo(() => validateOrg(normalizedSlug, normalizedName), [normalizedSlug, normalizedName]);
  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitted(true);
    if (errors.length > 0) return;
    const org = await createOrg({ slug: normalizedSlug, name: normalizedName }).unwrap();
    navigate(`/app/orgs/${org.id}/sites/new`);
  }
  return <main className="no-orgs-page"><section className="dashboard-panel"><EmptyState title="Welcome to go-go-host" body="Create your first organization, then add a site and upload a bundle." />
    {submitted && errors.length ? <ErrorCallout title="Fix organization details" error={errors.join('\n')} /> : null}
    {createState.error ? <ErrorCallout title="Unable to create organization" error={apiErrorMessage(createState.error)} /> : null}
    <form className="no-orgs-page__form" onSubmit={onSubmit} noValidate>
      <label htmlFor={slugId}>Organization slug</label><input id={slugId} value={slug} onChange={(event) => setSlug(event.target.value)} placeholder="demo" />
      <label htmlFor={nameId}>Organization name</label><input id={nameId} value={name} onChange={(event) => setName(event.target.value)} placeholder="Demo Org" />
      <button type="submit" data-part="btn" disabled={createState.isLoading}>{createState.isLoading ? 'Creating…' : 'Create organization'}</button>
    </form>
  </section></main>;
}
