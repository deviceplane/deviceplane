import React, { useEffect } from 'react';
import styled from 'styled-components';

import '../../lib/xterm.css';
import config from '../../config';
import Card from '../../components/card';
import { Column } from '../../components/core';

const Terminal = styled(Column)``;

Terminal.defaultProps = {
  padding: 3,
  border: 0,
  borderRadius: 1,
  height: 10,
};

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

const DeviceSsh = ({
  route: {
    data: { params, device },
  },
}) => {
  useEffect(() => {
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
          term.on('data', function(data) {
            stream.write(data);
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
        `${config.wsEndpoint}/projects/${params.project}/devices/${device.id}/ssh`,
        ['binary']
      ),
      username: '',
    });
  }, []);

  return (
    <Card size="xlarge">
      <Terminal bg="grays.0">
        <Column id="terminal" flex={1} />
      </Terminal>
    </Card>
  );
};

export default DeviceSsh;
