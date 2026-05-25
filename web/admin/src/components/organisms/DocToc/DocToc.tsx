import type { Heading } from '../../molecules/MarkdownRenderer';
import './DocToc.css';

export interface DocTocProps {
  headings: Heading[];
}

/**
 * DocToc renders a list of document headings as a clickable TOC
 * inside an OS1 window panel.
 */
export function DocToc({ headings }: DocTocProps) {
  if (headings.length === 0) return null;

  // Only show h2 and h3 in the TOC (skip the doc title h1)
  const tocItems = headings.filter((h) => h.level >= 2 && h.level <= 3);

  if (tocItems.length === 0) return null;

  const handleClick = (id: string) => {
    const el = document.getElementById(id);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  };

  return (
    <nav className="doc-toc" aria-label="Table of contents">
      <ul className="doc-toc__list">
        {tocItems.map((h) => (
          <li
            key={h.id}
            className={`doc-toc__item doc-toc__item--level-${h.level}`}
          >
            <button
              type="button"
              className="doc-toc__link"
              title={h.text}
              onClick={() => handleClick(h.id)}
            >
              {h.text}
            </button>
          </li>
        ))}
      </ul>
    </nav>
  );
}
