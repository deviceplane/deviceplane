import React, { useMemo, useState } from 'react';
import { useNavigation } from 'react-navi';

import api from '../../api';
import Card from '../../components/card';
import Table from '../../components/table';
import {
  DeviceLabelMulti,
  DeviceLabelKey,
} from '../../components/device-label';
import {
  Column,
  Row,
  Text,
  Checkbox,
  MultiSelect,
  Button,
  Icon,
  toaster,
} from '../../components/core';
import { labelColor } from '../../helpers/labels';

const supportedMetrics = [
  {
    description: `Provides the count and status of the devices in your project.`,
    name: 'devices',
    labels: [],
    properties: [],
  },
];

const Project = ({
  route: {
    data: { params, metrics, devices },
  },
}) => {
  const [editRow, setEditRow] = useState();
  const tableData = useMemo(
    () =>
      supportedMetrics.map(supportedMetric => {
        const metric = metrics.find(
          ({ name }) => name === supportedMetric.name
        );
        if (metric) {
          metric.properties = metric.properties || [];
          return {
            ...supportedMetric,
            ...metric,
            enabled: true,
          };
        }
        return {
          ...supportedMetric,
          enabled: false,
        };
      }),
    [metrics]
  );

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

  const navigation = useNavigation();

  const saveEdit = async () => {
    let newMetrics = tableData.map(metric => {
      if (metric.name === editRow.name) {
        return editRow;
      }
      return metric;
    });
    newMetrics = newMetrics.filter(({ enabled }) => !!enabled);

    try {
      await api.updateProjectMetricsConfig({
        projectId: params.project,
        data: newMetrics.map(({ name, properties, labels }) => ({
          name,
          properties,
          labels,
        })),
      });
      setEditRow(null);
      toaster.success('Metric updated.');
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

  const columns = useMemo(() => {
    return [
      {
        Header: 'Metric',
        accessor: 'name',
        Cell: ({ cell: { value } }) => (
          <Column>
            <Text fontSize={3} style={{ textTransform: 'capitalize' }}>
              {value}
            </Text>
            <Text fontSize={0} color="grays.8">{`deviceplane.${value}`}</Text>
          </Column>
        ),
      },
      {
        Header: 'Description',
        accessor: 'description',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ row: { original }, cell: { value } }) =>
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
          <Row title="When enabled, a Datadog tag with the device name is included.">
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
        Header: 'Enabled',
        accessor: 'enabled',
        Cell: ({ cell: { value }, row: { original } }) => {
          const editing = editRow && editRow.name === original.name;
          return (
            <Checkbox
              checked={editing ? editRow.enabled : value}
              onChange={enabled =>
                setEditRow({
                  ...(editing ? editRow : original),
                  enabled,
                })
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
            <Row justifyContent="flex-end">
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
            <Row justifyContent="flex-end">
              <Button
                title={<Icon icon="edit" size={16} color="primary" />}
                variant="icon"
                onClick={() => setEditRow(original)}
              />
            </Row>
          ),
        minWidth: '50px',
        maxWidth: '50px',
      },
    ];
  }, [editRow]);

  return (
    <Card
      title="Project Metrics"
      size="full"
      subtitle="These metrics provide insights over your entire project."
    >
      <Table columns={columns} data={tableData} editRow={editRow} />
    </Card>
  );
};

export default Project;
