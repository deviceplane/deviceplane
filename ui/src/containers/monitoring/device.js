import React, { useMemo, useState } from 'react';
import { useNavigation } from 'react-navi';

import api from '../../api';
import Card from '../../components/card';
import Table from '../../components/table';
import Popup from '../../components/popup';
import {
  DeviceLabelKey,
  DeviceLabelMulti,
} from '../../components/device-label';
import {
  Column,
  Row,
  Button,
  Text,
  Checkbox,
  MultiSelect,
  Icon,
  toaster,
} from '../../components/core';
import { getMetricLabel } from '../../helpers/metrics';
import DeviceMetricsForm from './device-metrics-form';
import { labelColor } from '../../helpers/labels';

const Device = ({
  route: {
    data: { params, devices, metrics },
  },
}) => {
  const [metricToDelete, setMetricToDelete] = useState();
  const [showMetricsForm, setShowMetricsForm] = useState();
  const [editRow, setEditRow] = useState();
  const navigation = useNavigation();

  const labelsOptions = useMemo(
    () =>
      [
        ...new Set(
          devices.reduce(
            (options, device) => [...options, ...Object.keys(device.labels)],
            []
          )
        ),
      ].map(
        label => ({
          label,
          value: label,
          props: {
            color: labelColor(label),
          },
        }),
        []
      ),
    [devices]
  );

  const hideMetricsForm = () => setShowMetricsForm(false);

  const submitDelete = async () => {
    setMetricToDelete(null);
    try {
      await api.updateDeviceMetricsConfig({
        projectId: params.project,
        data: metrics.filter(({ name }) => name !== metricToDelete.name),
      });
      toaster.success('Metric deleted.');
      navigation.refresh();
    } catch (error) {
      toaster.danger('Metric deletion failed.');
      console.error(error);
    }
  };

  const saveEdit = async () => {
    try {
      await api.updateDeviceMetricsConfig({
        projectId: params.project,
        data: metrics.map(metric =>
          metric.name === editRow.name ? editRow : metric
        ),
      });
      toaster.success('Metric updated.');
      setEditRow(null);
      navigation.refresh();
    } catch (error) {
      toaster.danger('Metric update failed.');
      console.error(error);
    }
  };

  const updateMetricProperty = (property, value, metric) => {
    setEditRow({
      ...metric,
      properties: value
        ? [...metric.properties, property]
        : metric.properties.filter(p => p !== property),
    });
  };

  const columns = useMemo(
    () => [
      {
        Header: 'Metric',
        accessor: 'name',
        Cell: ({ cell: { value } }) => (
          <Column>
            <Text fontSize={3}>{getMetricLabel(value)}</Text>
            <Text fontSize={0} color="grays.8">{`deviceplane.${value}`}</Text>
          </Column>
        ),
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value }, row: { original } }) =>
          editRow && editRow.name === original.name ? (
            <MultiSelect
              multi
              value={editRow.labels.map(label => ({
                label,
                value: label,
                props: { color: labelColor(label) },
              }))}
              options={labelsOptions}
              multiComponent={DeviceLabelMulti}
              onChange={(value, props) => {
                if (props.action === 'remove-value') {
                  setEditRow({
                    ...editRow,
                    labels: editRow.labels.filter(
                      label => label !== props.removedValue.value
                    ),
                  });
                } else {
                  setEditRow({
                    ...editRow,
                    labels: value.map(({ value }) => value),
                  });
                }
              }}
              placeholder="Select labels"
              noOptionsMessage={() => (
                <Text>
                  There are no <strong>Labels</strong>.
                </Text>
              )}
            />
          ) : (
            <Row
              onClick={() => setEditRow(original)}
              style={{ cursor: 'pointer' }}
            >
              {value.map(label => (
                <DeviceLabelKey key={label} label={label} />
              ))}
            </Row>
          ),
      },
      {
        id: 'device',
        accessor: ({ properties }) => properties.includes('device'),
        Header: (
          <Row
            alignItems="center"
            title="When enabled, a Datadog tag with the device name is included."
          >
            <Text marginLeft={1}>Device</Text>
          </Row>
        ),
        Cell: ({ cell: { value }, row: { original } }) => {
          const editing = editRow && editRow.name === original.name;
          return (
            <Checkbox
              checked={editing ? editRow.properties.includes('device') : value}
              onChange={v =>
                updateMetricProperty('device', v, editing ? editRow : original)
              }
            />
          );
        },
        minWidth: '100px',
        maxWidth: '100px',
        cellStyle: {
          justifyContent: 'center',
        },
      },
      {
        Header: ' ',
        Cell: ({ row: { original } }) =>
          editRow && editRow.name === original.name ? (
            <Row>
              <Button
                title={<Icon icon="floppy-disk" size={16} color="primary" />}
                variant="icon"
                onClick={saveEdit}
              />
              <Button
                title={<Icon icon="cross" size={16} color="white" />}
                variant="icon"
                onClick={() => setEditRow(null)}
                marginLeft={3}
              />
            </Row>
          ) : (
            <Row>
              <Button
                title={<Icon icon="edit" size={16} color="primary" />}
                variant="icon"
                onClick={() => setEditRow(original)}
              />
              <Button
                title={<Icon icon="trash" size={16} color="red" />}
                variant="icon"
                marginLeft={3}
                onClick={() => setMetricToDelete(original)}
              />
            </Row>
          ),
        minWidth: '80px',
        maxWidth: '80px',
        cellStyle: {
          justifyContent: 'flex-end',
        },
      },
    ],
    [editRow]
  );

  const tableData = useMemo(() => metrics, [metrics]);

  return (
    <>
      <Card
        title="Device Metrics"
        subtitle="These metrics provide information on the state of each device."
        size="full"
        actions={[
          {
            title: 'Add Device Metrics',
            onClick: () => setShowMetricsForm(true),
          },
        ]}
        maxHeight="100%"
      >
        <Table
          data={tableData}
          columns={columns}
          editRow={editRow}
          placeholder={
            <Text>
              There are no <strong>Device Metrics</strong>.
            </Text>
          }
        />
        <Popup show={showMetricsForm} onClose={hideMetricsForm}>
          <DeviceMetricsForm
            params={params}
            metrics={metrics}
            devices={devices}
            close={hideMetricsForm}
          />
        </Popup>
        <Popup show={!!metricToDelete} onClose={() => setMetricToDelete(null)}>
          <Card border title="Delete Device Metric" size="large">
            <Text>
              You are about to delete the{' '}
              <strong>
                {metricToDelete && getMetricLabel(metricToDelete.name)}
              </strong>{' '}
              metric.
            </Text>
            <Button
              marginTop={5}
              title="Delete"
              onClick={submitDelete}
              variant="danger"
            />
          </Card>
        </Popup>
      </Card>
    </>
  );
};

export default Device;
