export interface TimestampProps { value?: string; mode?: 'absolute' | 'relative'; }
export function Timestamp({ value, mode = 'absolute' }: TimestampProps) {
  if (!value || value.startsWith('0001-')) return <span title="No timestamp">—</span>;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return <span title={value}>Invalid date</span>;
  if (mode === 'relative') {
    const diff = Math.round((date.getTime() - Date.now()) / 1000);
    const abs = Math.abs(diff);
    const unit = abs > 3600 ? `${Math.round(abs / 3600)}h` : abs > 60 ? `${Math.round(abs / 60)}m` : `${abs}s`;
    return <time dateTime={date.toISOString()}>{diff < 0 ? `${unit} ago` : `in ${unit}`}</time>;
  }
  return <time dateTime={date.toISOString()}>{date.toLocaleString()}</time>;
}
