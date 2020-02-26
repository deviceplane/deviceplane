import React, { useEffect, useState, useRef, useMemo, Fragment } from 'react';
import styled from 'styled-components';
import { useTable, useSortBy, useRowSelect } from 'react-table';
import { Terminal as XTerm } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';

import '../lib/xterm.css';
import api from '../api';
import config from '../config';
import storage from '../storage';
import segment from '../lib/segment';
import Card from '../components/card';
import Popup from '../components/popup';
import Table, { SelectColumn } from '../components/table';
import { renderLabels } from '../helpers/labels';
import {
  Text,
  Button,
  Icon,
  Row,
  Column,
  Input,
  Select,
} from '../components/core';

var process = require('process');

// Hacks to get ssh2 working in browser
process.binding = function() {
  return {};
};
process.version = '1.0.0';

var Client = require('ssh2/lib/client');
var ws = require('websocket-stream');

const Container = styled(Column)``;

Container.defaultProps = {
  flex: 1,
  bg: 'black',
  overflow: 'hidden',
};

const Terminal = styled(Column)`
  position: absolute;
  overflow: hidden;
  flex: 1;
  width: 100%;
  height: 100%;
  visibility: ${props => (props.show ? 'visible' : 'hidden')};
`;

const Tabs = styled(Row)`
  height: 32px;
  flex: 1;
  overflow-x: scroll;
  scrollbar-width: none;
  -ms-overflow-style: none;

  ::-webkit-scrollbar {
    width: 0;
    height: 0;
    background: transparent;
  }
`;

const CloseButton = styled(Row)`
  cursor: pointer;
  border-radius: 50%;
  padding: 3px;
`;

const AddButton = styled(Row).attrs({ as: 'button' })`
  cursor: pointer;
  appearance: none;
  outline: none;
  border: none;
  background-color: ${props => props.theme.colors.grays[0]};
  border-left: 1px solid ${props => props.theme.colors.grays[8]};
  border-radius: 0;
  flex-shrink: 0;
  padding: 6px;
  &:hover {
    background-color: ${props => props.theme.colors.grays[3]};
  }
`;

const Tab = styled(Row).attrs({ as: 'button' })`
  position: relative;
  appearance: none;
  outline: none;
  border: none;
  border-radius: 0;
  cursor: pointer;
  font-size: 16px;
  font-weight: 500;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  flex: 1 0 0%;
  min-width: 150px;
  height: 32px;
  padding: 0 4px;
  border-right: 1px solid ${props => props.theme.colors.grays[8]};
  color: ${props =>
    props.active ? props.theme.colors.primary : props.theme.colors.white};
  background-color: ${props =>
    props.active ? props.theme.colors.grays[3] : props.theme.colors.grays[1]};
  & button {
    visibility: ${props => (props.active ? 'visible' : 'hidden')};
  }
  &:last-child {
    border-right: none;
  }
  &:hover ${CloseButton} {
    visibility: visible;
  }

  & ${CloseButton} {
    visibility: hidden;
    position: absolute;
    right: 8px;
    background-color: ${props =>
      props.active ? props.theme.colors.grays[3] : props.theme.colors.grays[1]};
  }
  & ${CloseButton}:hover {
    background-color: ${props =>
      props.active ? props.theme.colors.grays[5] : props.theme.colors.grays[4]};
  }
`;

const Session = ({ index, project, device, privateKey, show, onTitle }) => {
  const terminalNode = useRef();
  const client = useRef(new Client()).current;
  const term = useRef(
    new XTerm({ cursorStyle: 'bar', cursorWidth: '3', cursorBlink: true })
  ).current;
  const fitAddon = useRef(new FitAddon()).current;

  const startSSH = () => {
    const wndopts = { term: 'xterm' };

    client
      .on('ready', function() {
        client.shell(wndopts, function(err, stream) {
          if (err) {
            stream.end();
          }
          const {
            width,
            height,
          } = terminalNode.current.getBoundingClientRect();
          stream.setWindow(term.rows, term.cols, width, height);
          stream
            .on('close', function() {
              client.end();
            })
            .on('data', function(data) {
              term.write(data.toString());
            })
            .stderr.on('data', function(data) {
              term.write(data.toString());
            });
          term.onData(data => {
            stream.write(data);
          });
          term.onResize(({ rows, cols }) => {
            const {
              width,
              height,
            } = terminalNode.current.getBoundingClientRect();
            stream.setWindow(rows - 1, cols, width, height - 32);
          });
        });
      })
      .on('error', function(err) {
        term.write(err.toString());
      })
      .on('close', function() {
        term.write('\r\nConnection lost.\r\n');
      });

    const options = {
      sock: ws(
        `${config.wsEndpoint}/projects/${project}/devices/${device}/ssh`,
        ['binary']
      ),
      username: ' ',
    };

    if (privateKey) {
      options.privateKey = privateKey;
    }

    term.open(terminalNode.current);
    term.loadAddon(fitAddon);
    term.onTitleChange(title => onTitle(device, index, title));
    window.onresize = () => {
      fitAddon.fit();
    };
    fitAddon.fit();
    term.focus();

    client.connect(options);

    segment.track('Device SSH');
  };

  useEffect(() => {
    startSSH();
    return () => {
      client.end();
      fitAddon.dispose();
      term.dispose();
    };
  }, []);

  useEffect(() => {
    if (show) {
      term.focus();
    }
  }, [show]);

  return <Terminal ref={terminalNode} show={show} />;
};

const SessionTabs = ({ device, setActiveSession, deleteSession }) => {
  if (device && device.sessions) {
    return (
      <Tabs>
        {device.sessions.map(({ active, title }, i) => (
          <Tab key={i} active={active} onClick={() => setActiveSession(i)}>
            <Text color="inherit">{title || i + 1}</Text>
            <CloseButton
              marginLeft={2}
              onClick={e => {
                e.stopPropagation();
                deleteSession(i);
              }}
            >
              <Icon icon="cross" color="white" size={14} />
            </CloseButton>
          </Tab>
        ))}
      </Tabs>
    );
  }
  return null;
};

const SSH = ({
  route: {
    data: { params, devices: initialAllDevices },
    url: { query },
  },
}) => {
  const [allDevices, setAllDevices] = useState(initialAllDevices);
  const [devices, setDevices] = useState(
    query.devices
      ? query.devices.split(',').map((name, i) => ({
          name,
          active: i === 0,
          sessions: [{ active: true }],
        }))
      : []
  );
  const sshKeys = storage.get('sshKeys');
  const enableSSHKeys = storage.get('enableSSHKeys', params.project);
  const selectOptions = sshKeys
    ? sshKeys.map(({ name, key }) => ({ label: name, value: key }))
    : null;
  const [isKeyPopupVisible, setKeyPopupVisible] = useState(
    enableSSHKeys && selectOptions
  );
  const [isDevicesTableVisible, setDevicesTableVisible] = useState();
  const [searchInput, setSearchInput] = useState('');
  const [searchFocused, setSearchFocused] = useState();
  const [privateKey, setPrivateKey] = useState();

  useEffect(() => {
    setTimeout(() => {
      const intercomNode = document.querySelector('#intercom-container');
      if (intercomNode) {
        intercomNode.style.display = 'none';
        return () => {
          intercomNode.style.display = 'block';
        };
      }
    }, 500);
  }, []);

  const fetchDevices = async () => {
    try {
      const { data, headers } = await api.devices({
        projectId: params.project,
        queryString: `?search=${searchInput}`,
      });
      setAllDevices(data);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    if (devices.length === 0) {
      //window.close();
    } else {
      window.history.replaceState(
        null,
        null,
        `?devices=${devices.map(({ name }) => name).join(',')}`
      );
    }
  }, [devices]);

  useEffect(() => {
    fetchDevices();
  }, [searchInput]);

  const columns = useMemo(
    () => [
      SelectColumn,
      {
        Header: 'Name',
        accessor: 'name',
        minWidth: '200px',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value } }) =>
          value ? <Row marginBottom={-2}>{renderLabels(value)}</Row> : null,
        minWidth: '300px',
        maxWidth: '2fr',
      },
    ],
    []
  );

  const tableData = useMemo(
    () =>
      allDevices.filter(
        ({ status, name }) =>
          status === 'online' && !devices.find(device => device.name === name)
      ),
    [allDevices, devices]
  );

  const { selectedFlatRows, ...tableProps } = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy,
    useRowSelect
  );

  const deleteDevice = name =>
    setDevices(devices => {
      let wasDeviceActive = false;
      let index;
      const newDevices = devices.filter((device, i) => {
        if (device.name === name) {
          wasDeviceActive = device.active;
          index = i;
          return false;
        }
        return true;
      });
      if (wasDeviceActive) {
        if (newDevices[index]) {
          newDevices[index].active = true;
        } else if (newDevices[index - 1]) {
          newDevices[index - 1].active = true;
        }
      }
      return newDevices;
    });

  const addSession = () =>
    setDevices(devices =>
      devices.map(d =>
        d.active
          ? {
              ...d,
              sessions: [
                ...d.sessions.map(session => ({
                  ...session,
                  active: false,
                })),
                { active: true },
              ],
            }
          : d
      )
    );

  const setActiveSession = sessionIndex =>
    setDevices(devices =>
      devices.map(d =>
        d.active
          ? {
              ...d,
              sessions: d.sessions.map((session, i) => ({
                ...session,
                active: i === sessionIndex,
              })),
            }
          : d
      )
    );

  const deleteSession = sessionIndex =>
    setDevices(devices => {
      let deviceIndex;
      let deviceDeleted;
      const newDevices = devices
        .map((device, i) => {
          if (device.active) {
            deviceIndex = i;
            let wasSessionActive;
            const newSessions = device.sessions.filter((session, i) => {
              if (i === sessionIndex) {
                wasSessionActive = session.active;
                return false;
              }
              return true;
            });
            if (newSessions.length) {
              if (wasSessionActive) {
                if (newSessions[sessionIndex]) {
                  newSessions[sessionIndex].active = true;
                } else if (newSessions[sessionIndex - 1]) {
                  newSessions[sessionIndex - 1].active = true;
                }
              }
            } else {
              deviceDeleted = true;
            }
            return {
              ...device,
              sessions: newSessions,
            };
          }
          return device;
        })
        .filter(({ sessions }) => sessions.length);

      if (deviceDeleted) {
        if (newDevices[deviceIndex]) {
          newDevices[deviceIndex].active = true;
        } else if (newDevices[deviceIndex - 1]) {
          newDevices[deviceIndex - 1].active = true;
        }
      }

      return newDevices;
    });

  const setSessionTitle = (device, sessionIndex, title) => {
    if (title) {
      setDevices(devices =>
        devices.map(d =>
          d.name === device
            ? {
                ...d,
                sessions: d.sessions.map((session, i) =>
                  i === sessionIndex ? { ...session, title } : session
                ),
              }
            : d
        )
      );
    }
  };

  const addSelectedDevices = () => {
    setDevices(devices => [
      ...devices.map(device => ({ ...device, active: false })),
      ...selectedFlatRows.map(({ original: { name } }, i) => ({
        name,
        active: i === selectedFlatRows.length - 1,
        sessions: [{ active: true }],
      })),
    ]);
  };

  return (
    <>
      <Container>
        <Row borderBottom="1px solid" borderColor="grays.8">
          <Tabs>
            {devices.map(({ name, active }) => (
              <Tab
                key={name}
                active={active}
                onClick={() =>
                  setDevices(devices =>
                    devices.map(device => ({
                      ...device,
                      active: device.name === name,
                    }))
                  )
                }
              >
                <Text color="inherit">{name}</Text>
                {devices.length > 0 && (
                  <CloseButton
                    marginLeft={2}
                    onClick={e => {
                      e.stopPropagation();
                      deleteDevice(name);
                    }}
                  >
                    <Icon icon="cross" color="white" size={14} />
                  </CloseButton>
                )}
              </Tab>
            ))}
          </Tabs>

          <AddButton onClick={() => setDevicesTableVisible(true)}>
            <Icon icon="plus" color="primary" size={18} />
          </AddButton>
        </Row>
        <Row borderBottom="1px solid" borderColor="grays.8">
          <SessionTabs
            device={devices.find(({ active }) => active)}
            setActiveSession={setActiveSession}
            deleteSession={deleteSession}
          />

          {devices.length > 0 && (
            <AddButton onClick={addSession}>
              <Icon icon="plus" color="primary" size={18} />
            </AddButton>
          )}
        </Row>
        <Column
          width="100%"
          height="100%"
          padding={2}
          paddingBottom={3}
          position="relative"
          overflow="hidden"
        >
          {devices.map(({ active, name, sessions }) => (
            <Fragment key={name}>
              {sessions.map((session, i) => (
                <Session
                  key={i}
                  index={i}
                  show={active && session.active}
                  device={name}
                  project={params.project}
                  onTitle={setSessionTitle}
                />
              ))}
            </Fragment>
          ))}
        </Column>
      </Container>

      <Popup
        show={isDevicesTableVisible}
        onClose={() => {
          setDevicesTableVisible(false);
        }}
      >
        <Card
          border
          size="xlarge"
          title="Devices"
          right={
            <Row
              position="relative"
              alignItems="center"
              flex={1}
              minWidth="300px"
            >
              <Icon
                icon="search"
                size={16}
                color={searchFocused ? 'primary' : 'white'}
                style={{ position: 'absolute', left: 16 }}
              />
              <Input
                flex={1}
                bg="black"
                placeholder="Search devices by name or labels"
                paddingLeft={7}
                value={searchInput}
                onChange={e => setSearchInput(e.target.value)}
                onFocus={() => setSearchFocused(true)}
                onBlur={() => setSearchFocused(false)}
              />
            </Row>
          }
        >
          <Row marginBottom={3}>
            <Button
              title="SSH"
              variant="tertiary"
              disabled={selectedFlatRows.length === 0}
              onClick={() => {
                setDevicesTableVisible(false);
                addSelectedDevices();
              }}
            />
          </Row>
          <Table
            {...tableProps}
            placeholder={
              <Text>
                There are no eligible <strong>Devices</strong>.
              </Text>
            }
          />
        </Card>
      </Popup>

      <Popup
        show={isKeyPopupVisible}
        onClose={() => {
          setKeyPopupVisible(false);
          startSSH();
        }}
      >
        <Card border size="medium">
          <Select
            onChange={e => {
              setKeyPopupVisible(false);
              setPrivateKey(e.target.value);
            }}
            options={selectOptions}
            placeholder="Select a SSH key"
            none="There are no SSH keys"
          />
        </Card>
      </Popup>
    </>
  );
};

export default SSH;
