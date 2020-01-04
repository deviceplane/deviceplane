import React, { useLayoutEffect } from 'react';
import AceEditor from 'react-ace';
import 'brace/mode/yaml';
import 'brace/theme/pastel_on_dark';

const Editor = ({ readOnly, onChange, value, width, height }) => {
  useLayoutEffect(() => {
    if (readOnly) {
      document.querySelector('.ace_cursor').style.display = 'none';
      document.querySelector('.ace_gutter-active-line').style.display = 'none';
    }
  }, [readOnly]);

  return (
    <AceEditor
      fontSize={14}
      mode="yaml"
      showPrintMargin={false}
      width={width}
      height={height}
      tabSize={2}
      setOptions={{ showLineNumbers: true }}
      editorProps={{ $blockScrolling: Infinity }}
      readOnly={readOnly}
      highlightActiveLine={readOnly ? false : true}
      highlightGutterLine={readOnly ? false : true}
      value={value}
      onChange={onChange}
      theme="pastel_on_dark"
      style={{ borderRadius: '4px' }}
    />
  );
};

export default Editor;
