import React, { useState, useMemo } from 'react';
import styled from 'styled-components';
import { toaster, Icon } from 'evergreen-ui';

import theme from '../theme';
import Card from './card';
import Table from './table';
import Popup from './popup';
import { Text, Row, Button, Input } from './core';

const CellInput = styled(Input)`
  width: 100%;
  flex: 1;
`;

const EditableCell = ({ mode, value, autoFocus, onChange }) => {
  if (mode === 'edit' || mode === 'new') {
    return (
      <CellInput
        autoFocus={autoFocus}
        value={value}
        onChange={onChange}
        padding={1}
      />
    );
  }
  return (
    <Text textOverflow="ellipsis" overflow="hidden" whiteSpace="nowrap">
      {value}
    </Text>
  );
};

const EditableLabelTable = ({ data, onAdd, onRemove }) => {
  const [labelToRemove, setLabelToRemove] = useState();
  const [labels, setLabels] = useState(
    Object.keys(data)
      .map(key => ({
        key,
        value: data[key],
        editedKey: key,
        editedValue: data[key],
        mode: 'default',
      }))
      .sort((a, b) => {
        if (a.key < b.key) {
          return -1;
        }
        if (a.key > b.key) {
          return 1;
        }
        return 0;
      })
  );

  const columns = useMemo(
    () => [
      {
        Header: 'Key',
        Cell: ({ row: { index, original } }) => (
          <EditableCell
            mode={original.mode === 'edit' ? 'default' : original.mode}
            value={original.editedKey}
            onChange={e => editLabel(index, 'editedKey', e.target.value)}
            autoFocus
          />
        ),
        style: {
          flex: 3,
          minHeight: '56px',
          alignItems: 'center',
        },
      },
      {
        Header: 'Value',
        Cell: ({ row: { index, original } }) => (
          <EditableCell
            mode={original.mode}
            value={original.editedValue}
            onChange={e => editLabel(index, 'editedValue', e.target.value)}
          />
        ),
        style: {
          flex: 3,
          alignItems: 'center',
          minHeight: '56px',
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
                    !(original.editedKey && original.editedValue) ||
                    (original.mode === 'new' &&
                      data.find(({ key }) => key === original.editedKey))
                  }
                  onClick={() => saveLabel(original, index)}
                  variant="icon"
                />
                <Button
                  title={
                    <Icon
                      icon="cross"
                      size={16}
                      color={theme.colors.grays[5]}
                    />
                  }
                  variant="icon"
                  marginLeft={3}
                  onClick={() => cancelEdit(index, original.mode)}
                />
              </Row>
            );
          }
          return (
            <Row>
              <Button
                title={
                  <Icon icon="edit" size={16} color={theme.colors.primary} />
                }
                variant="icon"
                onClick={() => setEdit(index)}
              />
              <Button
                title={<Icon icon="trash" size={16} color={theme.colors.red} />}
                variant="icon"
                marginLeft={3}
                onClick={() => setLabelToRemove(original)}
              />
            </Row>
          );
        },
        style: {
          alignItems: 'center',
          justifyContent: 'flex-end',
          minHeight: '56px',
        },
      },
    ],
    []
  );
  const tableData = useMemo(() => labels, [labels]);

  const createLabel = (label = { key: '', value: '', mode: 'new' }) => {
    setLabels([...labels, label]);
  };

  const setEdit = index => {
    setLabels(labels =>
      labels.map((label, i) =>
        index === i ? { ...label, mode: 'edit' } : label
      )
    );
  };

  const cancelEdit = (index, mode) => {
    setLabels(labels =>
      mode === 'new'
        ? labels.filter((_, i) => i !== index)
        : labels.map((label, i) =>
            index === i
              ? {
                  ...label,
                  editedKey: label.key,
                  editedValue: label.value,
                  mode: 'default',
                }
              : label
          )
    );
  };

  const editLabel = (index, property, value) => {
    setLabels(labels =>
      labels.map((label, i) =>
        i === index ? { ...label, [property]: value } : label
      )
    );
  };

  const saveLabel = async (label, index) => {
    try {
      await onAdd({ key: label.editedKey, value: label.editedValue });

      setLabels(labels =>
        labels.map((label, i) =>
          i === index
            ? {
                ...label,
                key: label.editedKey,
                value: label.editedValue,
                mode: 'default',
              }
            : label
        )
      );
    } catch (error) {
      toaster.danger('Label was not saved.');
      console.log(error);
    }
  };

  const removeLabel = async () => {
    try {
      await onRemove(labelToRemove.key);
      setLabels(labels =>
        labels.filter(label => label.key !== labelToRemove.key)
      );
    } catch (error) {
      toaster.danger('Label was not removed.');
      console.log(error);
    }
    setLabelToRemove(null);
  };

  return (
    <>
      <Card
        title="Labels"
        size="xlarge"
        actions={[{ title: 'Add Label', onClick: () => createLabel() }]}
      >
        <Table
          columns={columns}
          data={tableData}
          placeholder={
            <Text>
              There are no <strong>Labels</strong>.
            </Text>
          }
        />
      </Card>
      <Popup show={!!labelToRemove} onClose={() => setLabelToRemove(null)}>
        <Card title="Remove Label" border>
          <Text>
            You are about to remove the{' '}
            <strong>{labelToRemove && labelToRemove.key}</strong> label.
          </Text>
          <Button
            marginTop={5}
            title="Remove"
            onClick={removeLabel}
            variant="danger"
          />
        </Card>
      </Popup>
    </>
  );
};

export default EditableLabelTable;
