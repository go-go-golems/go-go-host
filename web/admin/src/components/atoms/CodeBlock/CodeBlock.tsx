import './CodeBlock.css';
export interface CodeBlockProps { code: string; language?: string; }
export function CodeBlock({ code, language = 'text' }: CodeBlockProps) {
  return <pre className="code-block" data-language={language}><code>{code}</code></pre>;
}
