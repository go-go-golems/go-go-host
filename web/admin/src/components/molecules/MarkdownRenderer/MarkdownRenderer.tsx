import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import './MarkdownRenderer.css';

export interface MarkdownRendererProps {
  /** Raw markdown source string. */
  content: string;
  /** Optional extra class on the wrapper. */
  className?: string;
}

/**
 * MarkdownRenderer renders raw markdown into OS1-styled dashboard HTML.
 *
 * It uses react-markdown with remark-gfm for tables, strikethrough, etc.
 * The component owns its own CSS which normalizes headings, prose,
 * code blocks, tables, and lists into the OS1 dashboard font scale.
 */
export function MarkdownRenderer({ content, className }: MarkdownRendererProps) {
  return (
    <article className={`markdown-renderer ${className ?? ''}`}>
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
    </article>
  );
}
