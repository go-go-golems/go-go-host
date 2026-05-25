import type { ReactNode } from 'react';
import './EmptyState.css';
export interface EmptyStateProps { title: string; body?: string; action?: ReactNode; }
export function EmptyState({ title, body, action }: EmptyStateProps) {
  return <section className="empty-state"><h2>{title}</h2>{body ? <p>{body}</p> : null}{action ? <div className="empty-state__action">{action}</div> : null}</section>;
}
