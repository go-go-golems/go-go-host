import { Link, useParams } from 'react-router-dom';
import { docs, docsBySection } from '../../services/docs/docs-data';
import type { DocSection } from '../../services/docs/docs-data';
import './DocsIndexPage.css';

const sectionLabels: Record<DocSection, string> = {
  Tutorial: 'Tutorials',
  GeneralTopic: 'Reference',
  Example: 'Examples',
  Application: 'Application',
  Other: 'Other',
};

export function DocsIndexPage() {
  const { orgId } = useParams<{ orgId: string }>();
  const groups = docsBySection();

  return (
    <div className="docs-index-page">
      <section className="dashboard-panel docs-index-page__header">
        <header>
          <h1>Documentation</h1>
          <p>Learn how to build, deploy, and operate go-go-host apps and agents.</p>
        </header>
      </section>

      {(Object.keys(groups) as DocSection[]).map((section) => (
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
