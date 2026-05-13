import { useState, useCallback, useRef, type ReactNode } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import './MarkdownRenderer.css';
import './highlight-os1.css';

export interface MarkdownRendererProps {
  /** Raw markdown source string. */
  content: string;
  /** Optional extra class on the wrapper. */
  className?: string;
}

/**
 * A single code block with a copy button in the top-right corner.
 */
function CodeBlockWithCopy({ children, ...rest }: React.HTMLAttributes<HTMLPreElement> & { children?: ReactNode }) {
  const [copied, setCopied] = useState(false);
  const preRef = useRef<HTMLPreElement>(null);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleCopy = useCallback(() => {
    // Extract text from the <code> child
    const codeEl = preRef.current?.querySelector('code');
    const text = codeEl?.textContent ?? preRef.current?.textContent ?? '';
    void navigator.clipboard.writeText(text).then(() => {
      setCopied(true);
      clearTimeout(timerRef.current!);
      timerRef.current = setTimeout(() => setCopied(false), 1500);
    });
  }, []);

  return (
    <div className="markdown-renderer__code-wrap">
      <pre ref={preRef} {...rest}>
        {children}
      </pre>
      <button
        type="button"
        className="markdown-renderer__copy-btn"
        onClick={handleCopy}
        title={copied ? 'Copied!' : 'Copy to clipboard'}
        aria-label={copied ? 'Copied!' : 'Copy code to clipboard'}
      >
        {copied ? '✓' : '⧉'}
      </button>
    </div>
  );
}

/**
 * MarkdownRenderer renders raw markdown into OS1-styled dashboard HTML.
 *
 * It uses react-markdown with remark-gfm for tables, strikethrough, etc.
 * and rehype-highlight for syntax-highlighted code blocks.
 * Code blocks include a clipboard copy button.
 */
export function MarkdownRenderer({ content, className }: MarkdownRendererProps) {
  return (
    <article className={`markdown-renderer ${className ?? ''}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        rehypePlugins={[rehypeHighlight]}
        components={{ pre: CodeBlockWithCopy }}
      >
        {content}
      </ReactMarkdown>
    </article>
  );
}
