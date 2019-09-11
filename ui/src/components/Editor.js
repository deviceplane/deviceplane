import React, { Component } from 'react';
import './../App.css';
import AceEditor from 'react-ace';

export default class Editor extends Component {
  render() {
    const readOnly = this.props.hasOwnProperty('readOnly') && this.props.readOnly;
    const onChange = this.props.hasOwnProperty('onChange') ? this.props.onChange : null;
    return (
      <AceEditor
        // ref={this.editorRef}
        highlightActiveLine
        focus
        // annotations={getAnnotations(error)
        // markers={getMarkers(error)}
        fontSize={14}
        mode="yaml"
        theme="chrome"
        showPrintMargin={false}
        width={this.props.width}
        height={this.props.height}
        tabSize={2}
        setOptions={{ showLineNumbers: true }}
        editorProps={{ $blockScrolling: Infinity }}
        readOnly={readOnly}
        value={this.props.value}
        onChange={onChange}
        // onLoad={this.handleLoad}
      />
    );
  }
}