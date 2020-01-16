import React, { useState, useMemo, useEffect } from 'react';
import { toaster, Icon } from 'evergreen-ui';

import theme from '../theme';
import storage from '../storage';
import Card from '../components/card';
import Table from '../components/table';
import {
  Row,
  Column,
  Button,
  Text,
  Input,
  Textarea,
  Link,
} from '../components/core';

const EditableCell = ({
  mode,
  value,
  autoFocus,
  onChange,
  textArea,
  hideKey,
}) => {
  if (mode === 'edit' || mode === 'new') {
    const Component = textArea ? Textarea : Input;
    return (
      <Component
        autoFocus={autoFocus}
        value={value}
        onChange={onChange}
        padding={1}
        rows={20}
      />
    );
  }
  return <Text>{hideKey ? 'Edit to view key' : value}</Text>;
};

const SSHKeys = () => {
  const [sshKeys, setSSHKeys] = useState(
    (storage.get('sshKeys') || []).map(sshKey => ({
      ...sshKey,
      editedKey: sshKey.key,
      editedName: sshKey.name,
      mode: 'default',
    }))
  );

  useEffect(() => {
    storage.set('sshKeys', sshKeys);
  }, [sshKeys]);

  const addSSHKey = () =>
    setSSHKeys(keys => [
      { name: '', key: '', editedName: '', editedKey: '', mode: 'new' },
      ...keys,
    ]);

  const setMode = (index, mode) => {
    setSSHKeys(sshKeys =>
      sshKeys.map((sshKey, i) => (index === i ? { ...sshKey, mode } : sshKey))
    );
  };

  const cancelEdit = (index, mode) => {
    setSSHKeys(sshKeys =>
      mode === 'new'
        ? sshKeys.filter((_, i) => i !== index)
        : sshKeys.map((sshKey, i) =>
            index === i
              ? {
                  ...sshKey,
                  editedKey: sshKey.key,
                  editedName: sshKey.name,
                  mode: 'default',
                }
              : sshKey
          )
    );
  };

  const editSSHKey = (index, property, value) => {
    setSSHKeys(keys =>
      keys.map((sshKey, i) =>
        i === index ? { ...sshKey, [property]: value } : sshKey
      )
    );
  };

  const saveSSHKey = index => {
    try {
      setSSHKeys(sshKeys =>
        sshKeys.map((sshKey, i) =>
          i === index
            ? {
                ...sshKey,
                key: sshKey.editedKey,
                name: sshKey.editedName,
                mode: 'default',
              }
            : sshKey
        )
      );
      toaster.success('SSH key was saved to the browser successfully.');
    } catch (error) {
      toaster.danger('SSH key was not saved.');
      console.log(error);
    }
  };

  const deleteSSHKey = index => {
    try {
      setSSHKeys(sshKeys => sshKeys.filter((_, i) => i !== index));
      toaster.success('SSH key was successfully deleted.');
    } catch (error) {
      toaster.danger('SSH key was not deleted.');
      console.log(error);
    }
  };

  const columns = useMemo(
    () => [
      {
        Header: 'Name',
        accessor: 'name',
        Cell: ({ row: { original, index } }) => (
          <EditableCell
            autoFocus
            mode={original.mode}
            value={original.editedName}
            onChange={e => editSSHKey(index, 'editedName', e.target.value)}
          />
        ),
        style: {
          flex: '0 0 150px',
          alignSelf: 'flex-start',
        },
      },
      {
        Header: 'Key',
        accessor: 'key',
        Cell: ({ row: { original, index } }) => (
          <EditableCell
            hideKey
            textArea
            mode={original.mode}
            value={original.editedKey}
            onChange={e => editSSHKey(index, 'editedKey', e.target.value)}
          />
        ),
        style: {
          alignSelf: 'flex-start',
        },
      },
      {
        Header: ' ',
        Cell: ({ row: { index, original }, data }) => {
          if (original.mode === 'edit' || original.mode === 'new') {
            return (
              <Row>
                <Button
                  title={
                    <Icon
                      icon="floppy-disk"
                      size={16}
                      color={theme.colors.primary}
                    />
                  }
                  disabled={
                    !(original.editedKey && original.editedName) ||
                    (original.mode === 'new' &&
                      data.find(({ name }) => name === original.editedName))
                  }
                  onClick={() => saveSSHKey(index)}
                  variant="icon"
                />
                <Button
                  title={
                    <Icon icon="cross" size={16} color={theme.colors.white} />
                  }
                  variant="icon"
                  marginLeft={3}
                  onClick={() => cancelEdit(index, original.mode)}
                />
              </Row>
            );
          }
          if (original.mode === 'delete') {
            return (
              <>
                <Button
                  title={
                    <Icon
                      icon="tick-circle"
                      size={16}
                      color={theme.colors.red}
                    />
                  }
                  variant="icon"
                  onClick={() => deleteSSHKey(index)}
                />
                <Button
                  title={
                    <Icon icon="cross" size={16} color={theme.colors.white} />
                  }
                  variant="icon"
                  onClick={() => setMode(index, 'default')}
                  marginLeft={3}
                />
              </>
            );
          }
          return (
            <Row>
              <Button
                title={
                  <Icon icon="edit" size={16} color={theme.colors.primary} />
                }
                variant="icon"
                onClick={() => setMode(index, 'edit')}
              />
              <Button
                title={<Icon icon="trash" size={16} color={theme.colors.red} />}
                variant="icon"
                marginLeft={3}
                onClick={() => setMode(index, 'delete')}
              />
            </Row>
          );
        },
        style: {
          alignItems: 'center',
          justifyContent: 'flex-end',
          alignSelf: 'flex-start',
          flex: '0 0 100px',
        },
      },
    ],
    []
  );
  const tableData = useMemo(() => sshKeys, [sshKeys]);

  return (
    <>
      <Card
        border
        title="SSH Keys"
        subtitle={
          <Column>
            <Text>
              For security, these keys are only stored locally in the browser.
            </Text>
            <Text>
              <Link href="https://deviceplane.com/docs/device-variables/#authorized-ssh-keys">
                Learn more
              </Link>{' '}
              about configuring SSH keys.
            </Text>
          </Column>
        }
        size="xlarge"
        actions={[{ title: 'Add SSH Key', onClick: addSSHKey }]}
      >
        <Table
          columns={columns}
          data={tableData}
          placeholder={
            <Text>
              There are no <strong>SSH Keys</strong>.
            </Text>
          }
        />
      </Card>
    </>
  );
};

export default SSHKeys;
