import React, { useLayoutEffect, useState } from 'react';
import AceEditor from 'react-ace';
import 'brace/mode/yaml';
import 'brace/theme/pastel_on_dark';

import { Column } from './core';
import theme from '../theme';

const Editor = ({
  readOnly,
  onChange,
  value,
  width,
  height,
  maxLines = Infinity,
}) => {
  const [focused, setFocused] = useState();
  useLayoutEffect(() => {
    if (readOnly) {
      document.querySelector('.ace_cursor').style.display = 'none';
      document.querySelector('.ace_gutter-active-line').style.display = 'none';
    }
  }, [readOnly]);

  return (
    <Column
      flex={1}
      bg="grays.0"
      border={readOnly ? 'none' : 0}
      borderRadius={1}
      borderColor={!readOnly && focused ? 'primary' : 'white'}
      padding={2}
    >
      <AceEditor
        fontSize={14}
        mode="yaml"
        showPrintMargin={false}
        width={width}
        height={height}
        maxLines={maxLines}
        minLines={3}
        tabSize={2}
        setOptions={{ showLineNumbers: true }}
        editorProps={{ $blockScrolling: Infinity }}
        readOnly={readOnly}
        highlightActiveLine={readOnly ? false : true}
        highlightGutterLine={readOnly ? false : true}
        value={value}
        onChange={onChange}
        onFocus={() => setFocused(true)}
        onBlur={() => setFocused(false)}
        theme="pastel_on_dark"
      />
    </Column>
  );
};

export default Editor;
