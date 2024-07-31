// eslint-disable-next-line no-unused-vars
import React, { useRef, useEffect } from 'react';
import { EditorView } from '@codemirror/view';
import { EditorState } from '@codemirror/state';
import { basicSetup } from 'codemirror';
import { javascript } from '@codemirror/lang-javascript';
import { oneDark } from '@codemirror/theme-one-dark';

// Componente Editor
// eslint-disable-next-line react/prop-types
const Editor = ({ code, setCode }) => {
  const editorRef = useRef(null);

  useEffect(() => {
    if (editorRef.current) {
      const editorView = new EditorView({
        state: EditorState.create({
          doc: code,
          extensions: [
            basicSetup,
            javascript(),
            oneDark, // Aplica el tema como una extensión
            EditorView.updateListener.of(update => {
              if (update.docChanged) {
                setCode(update.state.doc.toString());
              }
            })
          ]
        }),
        parent: editorRef.current
      });

      return () => {
        editorView.destroy();
      };
    }
  }, [code, setCode]);

  return <div ref={editorRef} className="code-editor" />;
};

export default Editor;