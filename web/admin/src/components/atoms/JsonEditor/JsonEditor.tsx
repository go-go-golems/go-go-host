import { json } from '@codemirror/lang-json';
import { HighlightStyle, syntaxHighlighting } from '@codemirror/language';
import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands';
import { tags } from '@lezer/highlight';
import { EditorState } from '@codemirror/state';
import { EditorView, keymap, lineNumbers } from '@codemirror/view';
import { useEffect, useRef } from 'react';
import './JsonEditor.css';

const jsonHighlightStyle = HighlightStyle.define([
  { tag: tags.propertyName, color: '#004f8a', fontWeight: '700' },
  { tag: tags.string, color: '#0b6b2f' },
  { tag: tags.number, color: '#8a4b00' },
  { tag: tags.bool, color: '#7a2c91', fontWeight: '700' },
  { tag: tags.null, color: '#8a1111', fontWeight: '700' },
  { tag: tags.punctuation, color: 'var(--hc-color-fg)' },
]);

export interface JsonEditorProps {
  value: string;
  onChange: (value: string) => void;
  ariaLabel?: string;
}

export function JsonEditor({ value, onChange, ariaLabel = 'JSON editor' }: JsonEditorProps) {
  const hostRef = useRef<HTMLDivElement | null>(null);
  const viewRef = useRef<EditorView | null>(null);
  const onChangeRef = useRef(onChange);

  useEffect(() => { onChangeRef.current = onChange; }, [onChange]);

  useEffect(() => {
    if (!hostRef.current) return;
    const view = new EditorView({
      parent: hostRef.current,
      state: EditorState.create({
        doc: value,
        extensions: [
          lineNumbers(),
          history(),
          json(),
          syntaxHighlighting(jsonHighlightStyle),
          keymap.of([indentWithTab, ...defaultKeymap, ...historyKeymap]),
          EditorView.lineWrapping,
          EditorView.updateListener.of((update) => {
            if (update.docChanged) onChangeRef.current(update.state.doc.toString());
          }),
          EditorView.theme({
            '&': { minHeight: '8rem' },
            '.cm-scroller': { fontFamily: 'var(--hc-font-family)', minHeight: '8rem' },
          }),
        ],
      }),
    });
    view.dom.setAttribute('aria-label', ariaLabel);
    viewRef.current = view;
    return () => { view.destroy(); viewRef.current = null; };
  // Create the editor once for this component; external value sync is handled below.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    const view = viewRef.current;
    if (!view) return;
    const current = view.state.doc.toString();
    if (current === value) return;
    view.dispatch({ changes: { from: 0, to: current.length, insert: value } });
  }, [value]);

  return <div className="json-editor" data-part="json-editor" ref={hostRef} />;
}
