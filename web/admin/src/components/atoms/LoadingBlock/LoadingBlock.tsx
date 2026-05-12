import './LoadingBlock.css';
export function LoadingBlock({ lines = 3 }: { lines?: number }) {
  return <div className="loading-block" aria-label="Loading">{Array.from({ length: lines }, (_, i) => <span key={i} />)}</div>;
}
