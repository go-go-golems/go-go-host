import { Link, useParams } from 'react-router-dom';
import { useListDocsQuery } from '../../services/goGoHostApi';
import type { DocSection } from '../../services/types';
import { LoadingBlock } from '../../components/atoms/LoadingBlock';
import { ErrorCallout } from '../../components/atoms/ErrorCallout';
import './DocsIndexPage.css';

const sectionLabels: Record<string, string> = {
  Tutorial: 'Tutorials',
  GeneralTopic: 'Reference',
  Example: 'Examples',
  Application: 'Application',
  '': 'Other',
};

export function DocsIndexPage() {
  const { orgId } = useParams<{ orgId: string }>();
  const { data: docs, isLoading, error } = useListDocsQuery();

  if (isLoading) return <div className="docs-index-page"><LoadingBlock lines={6} /></div>;
  if (error) return <div className="docs-index-page"><ErrorCallout title="Failed to load docs" error={String(error)} /></div>;
  if (!docs || docs.length === 0) return <div className="docs-index-page"><ErrorCallout title="No docs available" error="No documentation entries were found." /></div>;

  // Group by section
  const groups: Record<string, typeof docs> = {};
  for (const doc of docs) {
    const section = doc.section || 'Other';
    (groups[section] ??= []).push(doc);
  }

  const sectionOrder: string[] = ['Tutorial', 'GeneralTopic', 'Example', 'Application', 'Other', ''];

  return (
    <div className="docs-index-page">
      <section className="dashboard-panel docs-index-page__header">
        <header>
          <h1>Documentation</h1>
          <p>Learn how to build, deploy, and operate go-go-host apps and agents.</p>
        </header>
      </section>

      {sectionOrder
        .filter((s) => groups[s]?.length)
        .map((section) => (
        <section className="dashboard-panel docs-index-page__section" key={section}>
          <header>
            <h2>{sectionLabels[section] ?? section}</h2>
          </header>
          <ul className="docs-index-page__list">
            {groups[section].map((doc) => (
              <li key={doc.slug} className="docs-index-page__item">
                <Link to={`/app/orgs/${orgId}/docs/${doc.slug}`} className="docs-index-page__link">
                  <span className="docs-index-page__title">{doc.title}</span>
                  <span className="docs-index-page__short">{doc.short}</span>
                  {doc.source === 'agent' && (
                    <span className="docs-index-page__badge" data-part="badge" data-tone="info">agent</span>
                  )}
                </Link>
              </li>
            ))}
          </ul>
        </section>
      ))}
    </div>
  );
}
