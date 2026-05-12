import { FormEvent, useMemo, useState } from 'react';
import { useOutletContext, useParams } from 'react-router-dom';
import { Checkbox } from '@go-go-golems/os-core';
import { CodeBlock, ErrorCallout, JsonEditor, LoadingBlock, StatusPill, Timestamp } from '../../components/atoms';
import { apiErrorMessage } from '../../services/errors';
import { useAddSiteDomainMutation, useDeleteSiteConfigMutation, useDeleteSiteDomainMutation, useGetSiteEnvironmentQuery, useListSiteCapabilitiesQuery, useListSiteConfigQuery, useListSiteDomainsQuery, useUpsertSiteCapabilityMutation, useUpsertSiteConfigMutation, useVerifySiteDomainMutation } from '../../services/goGoHostApi';
import type { Site } from '../../services/types';
import './SiteSettingsPage.css';

const safeCapabilities = ['express', 'ui.dsl', 'database', 'db', 'time', 'timer', 'assets'];

export function SiteSettingsPage() {
  const { siteId } = useParams();
  const { site } = useOutletContext<{ site: Site }>();
  const config = useListSiteConfigQuery(siteId ?? '', { skip: !siteId });
  const capabilities = useListSiteCapabilitiesQuery(siteId ?? '', { skip: !siteId });
  const domains = useListSiteDomainsQuery(siteId ?? '', { skip: !siteId });
  const environment = useGetSiteEnvironmentQuery(siteId ?? '', { skip: !siteId });
  const [upsertConfig, upsertConfigState] = useUpsertSiteConfigMutation();
  const [deleteConfig] = useDeleteSiteConfigMutation();
  const [upsertCapability, upsertCapabilityState] = useUpsertSiteCapabilityMutation();
  const [addDomain, addDomainState] = useAddSiteDomainMutation();
  const [verifyDomain] = useVerifySiteDomainMutation();
  const [deleteDomain] = useDeleteSiteDomainMutation();
  const [configKey, setConfigKey] = useState('theme.title');
  const [configValue, setConfigValue] = useState('{"text":"Hello"}');
  const [hostname, setHostname] = useState('www.example.com');
  const [configError, setConfigError] = useState('');
  const [capError, setCapError] = useState('');

  const capabilityRows = useMemo(() => {
    const byName = new Map((capabilities.data ?? []).map((cap) => [cap.capability, cap]));
    return safeCapabilities.map((name) => byName.get(name) ?? { siteId: site.id, capability: name, enabled: false, config: {}, updatedAt: '' });
  }, [capabilities.data, site.id]);

  async function submitConfig(event: FormEvent) {
    event.preventDefault();
    setConfigError('');
    try {
      const value = JSON.parse(configValue);
      await upsertConfig({ siteId: site.id, key: configKey, value }).unwrap();
    } catch (error) {
      setConfigError(error instanceof SyntaxError ? error.message : apiErrorMessage(error));
    }
  }

  async function toggleCapability(capability: string, enabled: boolean, configPayload: unknown) {
    setCapError('');
    try {
      await upsertCapability({ siteId: site.id, capability, enabled, config: configPayload ?? {} }).unwrap();
    } catch (error) {
      setCapError(apiErrorMessage(error));
    }
  }

  async function submitDomain(event: FormEvent) {
    event.preventDefault();
    if (!hostname.trim()) return;
    await addDomain({ siteId: site.id, hostname: hostname.trim() }).unwrap();
  }

  if (config.isLoading || capabilities.isLoading || domains.isLoading || environment.isLoading) return <section className="dashboard-panel"><LoadingBlock lines={8} /></section>;
  const firstError = config.error ?? capabilities.error ?? domains.error ?? environment.error;
  if (firstError) return <section className="dashboard-panel"><ErrorCallout title="Unable to load site settings" error={apiErrorMessage(firstError)} /></section>;

  return <div className="site-settings-page">
    <section className="dashboard-panel site-settings-page__intro">
      <header><h1>Site settings</h1><p>Configuration lives outside deployment bundles. Store <span className="site-settings-page__highlight site-settings-page__highlight--safe">non-secret values</span> here, keep <span className="site-settings-page__highlight site-settings-page__highlight--danger">secrets out of v1 runtime APIs</span>, and use domains to prepare traffic hosts.</p></header>
      <dl><div><dt>Primary host</dt><dd><code className="site-settings-page__identifier">{site.primaryHost}</code></dd></div><div><dt>Site ID</dt><dd><code className="site-settings-page__identifier">{site.id}</code></dd></div></dl>
    </section>

    <section className="dashboard-panel site-settings-page__section">
      <header><h2>Non-secret config</h2><p><span className="site-settings-page__highlight site-settings-page__highlight--info">JSON values</span> available to operators as site metadata. <span className="site-settings-page__highlight site-settings-page__highlight--warning">Secret storage is intentionally deferred.</span></p></header>
      <form className="site-settings-page__form" onSubmit={submitConfig}>
        <label>Key<input data-part="field-input" value={configKey} onChange={(e) => setConfigKey(e.target.value)} /></label>
        <label>JSON value<JsonEditor value={configValue} onChange={setConfigValue} ariaLabel="Site config JSON value" /></label>
        <button type="submit" data-part="btn" disabled={upsertConfigState.isLoading}>Save config</button>
      </form>
      {configError ? <ErrorCallout title="Config not saved" error={configError} /> : null}
      <table><thead><tr><th>Key</th><th>Value</th><th>Updated</th><th>Action</th></tr></thead><tbody>{(config.data ?? []).map((item) => <tr key={item.key}><td><code>{item.key}</code></td><td><CodeBlock code={JSON.stringify(item.value, null, 2)} /></td><td><Timestamp value={item.updatedAt} /></td><td><button type="button" data-part="btn" onClick={() => deleteConfig({ siteId: site.id, key: item.key })}>Delete</button></td></tr>)}</tbody></table>
    </section>

    <section className="dashboard-panel site-settings-page__section">
      <header><h2>Domains</h2><p>Add custom hostnames, copy the <span className="site-settings-page__highlight site-settings-page__highlight--info">verification token</span>, then mark them verified using the manual placeholder flow. <span className="site-settings-page__highlight site-settings-page__highlight--safe">Verified domains</span> are included on the next activation.</p></header>
      <form className="site-settings-page__form site-settings-page__form--inline" onSubmit={submitDomain}><label>Hostname<input data-part="field-input" value={hostname} onChange={(e) => setHostname(e.target.value)} /></label><button type="submit" data-part="btn" disabled={addDomainState.isLoading}>Add domain</button></form>
      {addDomainState.error ? <ErrorCallout title="Domain not added" error={apiErrorMessage(addDomainState.error)} /> : null}
      <table><thead><tr><th>Hostname</th><th>Status</th><th>Verification token</th><th>Verified</th><th>Actions</th></tr></thead><tbody>{(domains.data ?? []).map((domain) => <tr key={domain.id}><td><strong>{domain.hostname}</strong></td><td><StatusPill status={domain.status} tone={domain.status === 'verified' ? 'success' : 'warning'} /></td><td><code>{domain.verificationToken || '—'}</code></td><td><Timestamp value={domain.verifiedAt} /></td><td><button type="button" data-part="btn" onClick={() => verifyDomain({ siteId: site.id, domainId: domain.id })}>Verify</button> <button type="button" data-part="btn" onClick={() => deleteDomain({ siteId: site.id, domainId: domain.id })}>Delete</button></td></tr>)}</tbody></table>
    </section>

    <section className="dashboard-panel site-settings-page__section">
      <header><h2>Capabilities</h2><p>Org owners can toggle <span className="site-settings-page__highlight site-settings-page__highlight--info">site capability policy</span>. <code className="site-settings-page__highlight site-settings-page__highlight--danger">exec</code> and unrestricted <code className="site-settings-page__highlight site-settings-page__highlight--danger">fs</code> remain unavailable in hosted v1.</p></header>
      {capError ? <ErrorCallout title="Capability not updated" error={capError} /> : null}
      <table><thead><tr><th>Capability</th><th>Policy</th><th>Status</th><th>Config</th></tr></thead><tbody>{capabilityRows.map((cap) => <tr key={cap.capability}><td><code className={cap.enabled ? 'site-settings-page__highlight site-settings-page__highlight--safe' : undefined}>{cap.capability}</code></td><td><Checkbox label={cap.enabled ? 'Enabled' : 'Disabled'} checked={cap.enabled} disabled={upsertCapabilityState.isLoading} onChange={() => { void toggleCapability(cap.capability, !cap.enabled, cap.config); }} /></td><td><StatusPill status={cap.enabled ? 'enabled' : 'disabled'} tone={cap.enabled ? 'success' : 'danger'} /></td><td><CodeBlock code={JSON.stringify(cap.config ?? {}, null, 2)} /></td></tr>)}</tbody></table>
    </section>

    <section className="dashboard-panel site-settings-page__section">
      <header><h2>Environment and secrets</h2><p><span className="site-settings-page__highlight site-settings-page__highlight--warning">{environment.data?.message}</span></p></header>
      <div className="site-settings-page__env"><div><h3>Supported now</h3><ul>{environment.data?.supported.map((item) => <li className="site-settings-page__highlight site-settings-page__highlight--safe" key={item}>{item}</li>)}</ul></div><div><h3>Not supported in v1</h3><ul>{environment.data?.notSupported.map((item) => <li className="site-settings-page__highlight site-settings-page__highlight--danger" key={item}>{item}</li>)}</ul></div></div>
    </section>
  </div>;
}
