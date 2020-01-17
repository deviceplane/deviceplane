import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { Tooltip, Icon } from 'evergreen-ui';

import '../../lib/xterm.css';
import theme from '../../theme';
import config from '../../config';
import storage from '../../storage';
import Card from '../../components/card';
import Popup from '../../components/popup';
import { Row, Column, Select, Link } from '../../components/core';

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

const Terminal = styled(Column)``;

Terminal.defaultProps = {
  padding: 3,
  border: 0,
  borderRadius: 1,
  flex: 1,
};

const DeviceSsh = ({
  route: {
    data: { params, device },
  },
}) => {
  const sshKeys = storage.get('sshKeys');
  const enableSSHKeys = storage.get('enableSSHKeys', params.project);
  const selectOptions = sshKeys
    ? sshKeys.map(({ name, key }) => ({ label: name, value: key }))
    : null;
  const [showKeyPopup, setShowKeyPopup] = useState(
    enableSSHKeys && selectOptions
  );

  const startSSH = privateKey => {
    const conn = new Client();
    const term = new xterm();

    window.term = term;

    const wndopts = { term: 'xterm' };

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

    term.open(document.getElementById('terminal'), true);
    term.fit();
    term.clear();

    window.onresize = term.fit.bind(term);

    const options = {
      sock: ws(
        `${config.wsEndpoint}/projects/${params.project}/devices/${device.id}/ssh`,
        ['binary']
      ),
      username: ' ',
    };

    if (privateKey) {
      options.privateKey = privateKey;
    }

    conn.connect(options);
  };

  useEffect(() => {
    if (!showKeyPopup) {
      startSSH();
    }
  }, []);

  return (
    <>
      <Card size="full" height="100%">
        <Terminal bg="grays.0">
          <Column id="terminal" flex={1} />
        </Terminal>
      </Card>

      <Popup
        show={showKeyPopup}
        onClose={() => {
          setShowKeyPopup(false);
          startSSH();
        }}
      >
        <Card border>
          <Select
            onChange={({ value }) => {
              setShowKeyPopup(false);
              startSSH(value);
            }}
            options={selectOptions}
            placeholder="Select a SSH Key"
          />
        </Card>
      </Popup>
    </>
  );
};

export default DeviceSsh;
