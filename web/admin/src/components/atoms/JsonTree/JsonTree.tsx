import { CodeBlock } from '../CodeBlock';
export function JsonTree({ value }: { value: unknown }) {
  return <CodeBlock language="json" code={typeof value === 'string' ? value : JSON.stringify(value, null, 2)} />;
}
