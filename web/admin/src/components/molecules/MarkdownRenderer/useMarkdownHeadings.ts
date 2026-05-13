import { useState, useEffect } from 'react';

export interface Heading {
  id: string;
  level: number;
  text: string;
}

/**
 * Extracts headings from the rendered markdown DOM.
 *
 * Call this inside a component that renders MarkdownRenderer.
 * After the first paint the headings will be scanned from the article
 * element and returned as a stable array.
 */
export function useMarkdownHeadings(
  articleRef: React.RefObject<HTMLElement | null>,
  content: string,
): Heading[] {
  const [headings, setHeadings] = useState<Heading[]>([]);

  useEffect(() => {
    const el = articleRef.current;
    if (!el) return;

    const found: Heading[] = [];
    for (const h of el.querySelectorAll('h1, h2, h3, h4, h5, h6')) {
      const level = parseInt(h.tagName[1], 10);
      const text = h.textContent?.trim() ?? '';
      // rehype-highlight may have already assigned ids, otherwise derive one
      const id = h.id || text.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
      if (!h.id) h.id = id;
      found.push({ id, level, text });
    }
    setHeadings(found);
  }, [articleRef, content]);

  return headings;
}
