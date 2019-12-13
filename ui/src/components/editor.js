import React from 'react';
import AceEditor from 'react-ace';
import 'brace/mode/yaml';
import 'brace/theme/pastel_on_dark';

const Editor = ({ readOnly, onChange, value, width, height }) => (
  <AceEditor
    // ref={this.editorRef}
    highlightActiveLine
    // annotations={getAnnotations(error)
    // markers={getMarkers(error)}
    fontSize={14}
    mode="yaml"
    showPrintMargin={false}
    width={width}
    height={height}
    tabSize={2}
    setOptions={{ showLineNumbers: true }}
    editorProps={{ $blockScrolling: Infinity }}
    readOnly={readOnly}
    value={value}
    onChange={onChange}
    theme="pastel_on_dark"
    style={{ borderRadius: '4px' }}
    // onLoad={this.handleLoad}
  />
);

export default Editor;
