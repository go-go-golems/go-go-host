import { Link, useParams } from 'react-router-dom';
import { docBySlug, docs } from '../../services/docs/docs-data';
import { MarkdownRenderer } from '../../components/molecules/MarkdownRenderer';
import { EmptyState } from '../../components/atoms/EmptyState';
import './DocViewPage.css';

export function DocViewPage() {
  const { orgId, slug } = useParams<{ orgId: string; slug: string }>();
  const doc = slug ? docBySlug(slug) : undefined;

  if (!doc) {
    return (
      <div className="doc-view-page">
        <section className="dashboard-panel">
          <EmptyState title="Document not found" body={`No document with slug "${slug}" exists.`} />
          <Link to={`/app/orgs/${orgId}/docs`}>Back to documentation index</Link>
        </section>
      </div>
    );
  }

  return (
    <div className="doc-view-page">
      <section className="dashboard-panel doc-view-page__header">
        <nav className="doc-view-page__breadcrumb">
          <Link to={`/app/orgs/${orgId}/docs`}>Documentation</Link>
          <span className="doc-view-page__sep">›</span>
          <span>{doc.title}</span>
        </nav>
        <header>
          <h1>{doc.title}</h1>
          {doc.short && <p>{doc.short}</p>}
        </header>
      </section>

      <section className="dashboard-panel doc-view-page__body">
        <MarkdownRenderer content={doc.body} />
      </section>

      <section className="dashboard-panel doc-view-page__nav">
        <DocNavLinks currentSlug={doc.slug} orgId={orgId!} />
      </section>
    </div>
  );
}

function DocNavLinks({ currentSlug, orgId }: { currentSlug: string; orgId: string }) {
  const idx = docs.findIndex((d) => d.slug === currentSlug);
  const prev = idx > 0 ? docs[idx - 1] : null;
  const next = idx < docs.length - 1 ? docs[idx + 1] : null;

  return (
    <nav className="doc-view-page__prevnext">
      {prev ? (
        <Link to={`/app/orgs/${orgId}/docs/${prev.slug}`} className="doc-view-page__navlink">
          ‹ {prev.title}
        </Link>
      ) : <span />}
      {next ? (
        <Link to={`/app/orgs/${orgId}/docs/${next.slug}`} className="doc-view-page__navlink">
          {next.title} ›
        </Link>
      ) : <span />}
    </nav>
  );
}
