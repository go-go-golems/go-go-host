import './ConfirmActionDialog.css';

export interface ConfirmActionDialogProps {
  open: boolean;
  title: string;
  body: string;
  confirmLabel?: string;
  cancelLabel?: string;
  busy?: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmActionDialog({ open, title, body, confirmLabel = 'Confirm', cancelLabel = 'Cancel', busy = false, onConfirm, onCancel }: ConfirmActionDialogProps) {
  if (!open) return null;
  return <div className="confirm-action-dialog" role="presentation" onKeyDown={(event) => { if (event.key === 'Escape' && !busy) onCancel(); }}><section className="confirm-action-dialog__window" role="dialog" aria-modal="true" aria-labelledby="confirm-action-title"><h2 id="confirm-action-title">{title}</h2><p>{body}</p><footer><button type="button" data-part="btn" onClick={onCancel} disabled={busy} autoFocus>{cancelLabel}</button><button type="button" data-part="btn" onClick={onConfirm} disabled={busy}>{busy ? 'Working…' : confirmLabel}</button></footer></section></div>;
}
