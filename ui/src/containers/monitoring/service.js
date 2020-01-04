import React, { useMemo, useState, useEffect } from 'react';
import { useNavigation } from 'react-navi';
import { Tooltip, Icon, toaster } from 'evergreen-ui';

import storage from '../../storage';
import theme, { labelColors } from '../../theme';
import api from '../../api';
import { buildLabelColorMap } from '../../helpers/labels';
import Card from '../../components/card';
import Table from '../../components/table';
import Popup from '../../components/popup';
import {
  DeviceLabelKey,
  DeviceLabelMulti,
} from '../../components/device-label';
import { Button, Row, Text, Select, Checkbox } from '../../components/core';
import ServiceMetricsForm from './service-metrics-form';

const Service = ({
  route: {
    data: { params, applications, metrics, devices },
  },
}) => {
  const [selection, setSelection] = useState(
    storage.get('selectedService', params.project)
  );
  const [labelColorMap] = useState(
    buildLabelColorMap({}, labelColors, devices)
  );
  const [showMetricsForm, setShowMetricsForm] = useState();
  const [metricToDelete, setMetricToDelete] = useState();
  const [editRow, setEditRow] = useState();

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
            color: labelColorMap[label],
          },
        }),
        []
      ),
    [devices]
  );

  const navigation = useNavigation();

  const hideMetricsForm = () => setShowMetricsForm(false);
  const clearMetricToDelete = () => setMetricToDelete(null);

  let selectedMetrics = [];

  if (selection) {
    const serviceMetrics = metrics.find(
      ({ applicationId, service }) =>
        applicationId === selection.application.id &&
        service === selection.service
    );
    if (serviceMetrics) {
      selectedMetrics = serviceMetrics.exposedMetrics;
    }
  }

  const saveEdit = async () => {
    try {
      await api.updateServiceMetricsConfig({
        projectId: params.project,
        data: metrics.map(m => {
          if (m.service === selection.service) {
            return {
              ...m,
              exposedMetrics: m.exposedMetrics.map(em => {
                if (em.name === editRow.name) {
                  return editRow;
                }
                return em;
              }),
            };
          }
          return m;
        }),
      });
      toaster.success('Metric successfully updated.');
      setEditRow(null);
      navigation.refresh();
    } catch (error) {
      toaster.danger('Metric was not updated.');
      console.log(error);
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

  const submitDelete = async () => {
    clearMetricToDelete();
    try {
      await api.updateServiceMetricsConfig({
        projectId: params.project,
        data: metrics.map(m => {
          if (m.service === selection.service) {
            return {
              ...m,
              exposedMetrics: m.exposedMetrics.filter(
                ({ name }) => name !== metricToDelete.name
              ),
            };
          }
          return m;
        }),
      });
      toaster.success('Metric successfully deleted.');
      navigation.refresh();
    } catch (e) {
      console.log(e);
      toaster.danger('Metric was not deleted.');
    }
  };

  useEffect(() => {
    storage.set('selectedService', selection, params.project);
  }, [selection]);

  const tableData = useMemo(() => selectedMetrics, [selectedMetrics]);

  const columns = useMemo(
    () => [
      {
        Header: 'Metric',
        accessor: 'name',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ cell: { value }, row: { original } }) =>
          editRow && editRow.name === original.name ? (
            <Select
              multi
              options={labelsOptions}
              value={editRow.labels.map(label => ({
                label,
                value: label,
                props: { color: labelColorMap[label] },
              }))}
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
                <DeviceLabelKey
                  key={label}
                  label={label}
                  color={labelColorMap[label]}
                />
              ))}
            </Row>
          ),
      },
      {
        id: 'device',
        accessor: ({ properties }) => properties.includes('device'),
        Header: (
          <Row alignItems="center">
            <Tooltip content="When enabled, a Datadog tag with the device name is included.">
              <Icon icon="info-sign" size={10} color={theme.colors.primary} />
            </Tooltip>
            <Text marginLeft={1}>Device</Text>
          </Row>
        ),
        Cell: ({ cell: { value }, row: { original } }) => {
          const editing = editRow && editRow.name === original.name;
          return (
            <Checkbox
              value={editing ? editRow.properties.includes('device') : value}
              onChange={v =>
                updateMetricProperty('device', v, editing ? editRow : original)
              }
            />
          );
        },
        style: { flex: '0 0 125px', justifyContent: 'center' },
      },
      {
        Header: ' ',
        Cell: ({ row: { original } }) =>
          editRow && editRow.name === original.name ? (
            <Row>
              <Button
                title={
                  <Icon
                    icon="floppy-disk"
                    size={16}
                    color={theme.colors.primary}
                  />
                }
                variant="icon"
                onClick={saveEdit}
              />
              <Button
                title={
                  <Icon icon="cross" size={16} color={theme.colors.grays[5]} />
                }
                variant="icon"
                onClick={() => setEditRow(null)}
                marginLeft={3}
              />
            </Row>
          ) : (
            <Row>
              <Button
                title={
                  <Icon icon="edit" size={16} color={theme.colors.primary} />
                }
                variant="icon"
                onClick={() => setEditRow(original)}
              />
              <Button
                title={<Icon icon="trash" size={16} color={theme.colors.red} />}
                variant="icon"
                marginLeft={3}
                onClick={() => setMetricToDelete(original)}
              />
            </Row>
          ),
        style: {
          flex: '0 0 100px',
          justifyContent: 'flex-end',
        },
      },
    ],
    [editRow]
  );
  const selectOptions = useMemo(
    () =>
      applications
        .reduce((list, application) => {
          if (application.latestRelease) {
            return [
              ...list,
              ...Object.keys(application.latestRelease.config).map(service => ({
                application,
                service,
              })),
            ];
          }
          return list;
        }, [])
        .map(({ application, service }) => ({
          label: `${application.name}/${service}`,
          value: `${application.name}/${service}`,
          application,
          service,
        })),
    [applications]
  );

  return (
    <>
      <Row marginBottom={4} width={9}>
        <Select
          variant="black"
          onChange={setSelection}
          value={selection}
          options={selectOptions}
          placeholder="Select a Service"
          noOptionsMessage={() => (
            <Text>
              There are no <strong>Services</strong>.
            </Text>
          )}
        />
      </Row>
      <Card
        title="Service Metrics"
        subtitle="These are custom metrics you define on your services."
        size="full"
        actions={[
          {
            title: 'Add Service Metrics',
            onClick: () => setShowMetricsForm(true),
          },
        ]}
        disabled={!(selection && selection.service)}
        maxHeight="100%"
      >
        <Table
          data={tableData}
          columns={columns}
          placeholder={
            <Text>
              There are no <strong>Service Metrics</strong>.
            </Text>
          }
        />
      </Card>
      <Popup
        show={!!showMetricsForm}
        onClose={hideMetricsForm}
        overflow="visible"
      >
        <ServiceMetricsForm
          params={params}
          allMetrics={metrics}
          metrics={selectedMetrics}
          devices={devices}
          application={selection && selection.application}
          service={selection && selection.service}
          close={hideMetricsForm}
          labelColorMap={labelColorMap}
        />
      </Popup>
      <Popup show={!!metricToDelete} onClose={clearMetricToDelete}>
        <Card border title="Delete Service Metric">
          <Text>
            You are about to delete the{' '}
            <strong>{metricToDelete && metricToDelete.name}</strong> metric.
          </Text>
          <Button
            marginTop={5}
            title="Delete"
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </>
  );
};

export default Service;
