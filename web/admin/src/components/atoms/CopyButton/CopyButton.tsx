import { useState } from 'react';
import './CopyButton.css';
export interface CopyButtonProps { value: string; label?: string; }
export function CopyButton({ value, label = 'Copy' }: CopyButtonProps) {
  const [state, setState] = useState<'idle' | 'copied' | 'error'>('idle');
  async function copy() {
    try { await navigator.clipboard.writeText(value); setState('copied'); setTimeout(() => setState('idle'), 1200); }
    catch { setState('error'); }
  }
  return <button className="copy-button" data-part="btn" data-state={state} type="button" onClick={copy}>{state === 'copied' ? 'Copied' : state === 'error' ? 'Copy failed' : label}</button>;
}
