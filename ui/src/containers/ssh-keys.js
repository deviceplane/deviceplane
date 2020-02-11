import React, { useState, useMemo } from 'react';

import storage from '../storage';
import Card from '../components/card';
import Table from '../components/table';
import Field from '../components/field';
import {
  Row,
  Column,
  Button,
  Text,
  Link,
  Icon,
  toaster,
} from '../components/core';

const PrivateKeyRegexp = /^-----BEGIN RSA PRIVATE KEY-----(?:\r\n|\n)((?:[^:]+:\s*[\S].*(?:\r\n|\n))*)([\s\S]+)(?:\r\n|\n)-----END RSA PRIVATE KEY-----$/;

const EditableCell = ({
  mode,
  value,
  autoFocus,
  onChange,
  type,
  hideKey,
  error,
}) => {
  if (mode === 'edit' || mode === 'new') {
    return (
      <Field
        type={type}
        autoFocus={autoFocus}
        value={value}
        onChange={onChange}
        errors={error}
        rows={16}
        padding={1}
        marginBottom={0}
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
    setSSHKeys(sshKeys => {
      const validKey = PrivateKeyRegexp.test(sshKeys[index].editedKey);
      if (validKey) {
        const newSSHKeys = sshKeys.map((sshKey, i) =>
          i === index
            ? {
                ...sshKey,
                key: sshKey.editedKey,
                name: sshKey.editedName,
                mode: 'default',
              }
            : sshKey
        );
        storage.set('sshKeys', newSSHKeys);
        toaster.success('SSH key saved to the browser.');
        return newSSHKeys;
      } else {
        return sshKeys.map((sshKey, i) =>
          i === index
            ? {
                ...sshKey,
                error: { message: 'Invalid SSH Key, must be in RSA format.' },
              }
            : sshKey
        );
      }
    });
  };

  const deleteSSHKey = index => {
    try {
      setSSHKeys(sshKeys => sshKeys.filter((_, i) => i !== index));
      toaster.success('SSH key deleted.');
    } catch (error) {
      toaster.danger('SSH key deletion failed.');
      console.error(error);
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
      },
      {
        Header: 'Key',
        accessor: 'key',
        Cell: ({ row: { original, index } }) => (
          <EditableCell
            hideKey
            type="textarea"
            mode={original.mode}
            value={original.editedKey}
            error={original.error}
            onChange={e => editSSHKey(index, 'editedKey', e.target.value)}
          />
        ),
      },
      {
        Header: ' ',
        Cell: ({ row: { index, original }, data }) => {
          if (original.mode === 'edit' || original.mode === 'new') {
            return (
              <Row>
                <Button
                  title={<Icon icon="floppy-disk" size={16} color="primary" />}
                  disabled={
                    !(original.editedKey && original.editedName) ||
                    (original.mode === 'new' &&
                      data.find(({ name }) => name === original.editedName))
                  }
                  onClick={() => saveSSHKey(index)}
                  variant="icon"
                />
                <Button
                  title={<Icon icon="cross" size={16} color="white" />}
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
                  title={<Icon icon="tick-circle" size={16} color="red" />}
                  variant="icon"
                  onClick={() => deleteSSHKey(index)}
                />
                <Button
                  title={<Icon icon="cross" size={16} color="white" />}
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
                title={<Icon icon="edit" size={16} color="primary" />}
                variant="icon"
                onClick={() => setMode(index, 'edit')}
              />
              <Button
                title={<Icon icon="trash" size={16} color="red" />}
                variant="icon"
                marginLeft={3}
                onClick={() => setMode(index, 'delete')}
              />
            </Row>
          );
        },
        maxWidth: '50px',
        cellStyle: {
          alignItems: 'center',
          justifyContent: 'flex-end',
          alignSelf: 'flex-start',
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
              Keys must be in RSA format.{' '}
              <Link href="https://deviceplane.com/docs/variables/#authorized-ssh-keys">
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
