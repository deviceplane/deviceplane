import React, { Component } from 'react';
import { Pane, majorScale } from 'evergreen-ui';

import '../../vendor/xterm.css';
import config from '../../config.js';
import segment from '../../segment.js';

var process = require('process');

// Hacks to get ssh2 working in browser
process.binding = function() {
  return {};
};
process.version = '1.0.0';

var Client = require('ssh2/lib/client');
var ws = require('websocket-stream');
var xterm = require('xterm');
require('xterm/lib/addons/fit/fit');

export default class DeviceSsh extends Component {
  componentDidMount() {
    segment.page();

    var conn = new Client();
    var term = new xterm();
    var wndopts = { term: 'xterm' };

    window.term = term;

    // Store current size for initialization
    term.on('resize', function(rev) {
      wndopts.rows = rev.rows;
      wndopts.cols = rev.cols;
    });

    term.on('title', function(title) {
      document.title = title;
    });

    conn
      .on('ready', function() {
        conn.shell(wndopts, function(err, stream) {
          if (err) throw err;
          stream
            .on('close', function() {
              conn.end();
            })
            .on('data', function(data) {
              term.write(data.toString());
            })
            .stderr.on('data', function(data) {
              term.write(data.toString());
            });
          term.on('key', function(key, kev) {
            stream.write(key);
          });
          term.on('resize', function(rev) {
            stream.setWindow(rev.rows, rev.cols, 480, 640);
          });
        });
      })
      .on('error', function(err) {
        term.write(err.toString());
      })
      .on('close', function() {
        term.write('\r\nConnection lost.\r\n');
      });

    term.open(document.getElementById('terminal'));
    term.fit();
    window.onresize = term.fit.bind(term);

    conn.connect({
      sock: ws(
        `${config.wsEndpoint}/projects/${this.props.projectName}/devices/${this.props.device.id}/ssh`,
        ['binary']
      ),
      username: '',
      readyTimeout: 60000,
    });
  }

  render() {
    return <Pane id="terminal" margin={majorScale(4)} height="100%" />;
  }
}
