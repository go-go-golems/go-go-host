import './MetricCard.css';
export interface MetricCardProps { label: string; value: string | number; detail?: string; tone?: 'default' | 'danger' | 'success'; }
export function MetricCard({ label, value, detail, tone = 'default' }: MetricCardProps) { return <article className="metric-card" data-tone={tone}><span>{label}</span><strong>{value}</strong>{detail ? <small>{detail}</small> : null}</article>; }
