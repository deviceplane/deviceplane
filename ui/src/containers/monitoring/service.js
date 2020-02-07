import React, { useMemo, useState, useEffect } from 'react';
import { useNavigation } from 'react-navi';

import storage from '../../storage';
import api from '../../api';
import Card from '../../components/card';
import Table from '../../components/table';
import Popup from '../../components/popup';
import {
  DeviceLabelKey,
  DeviceLabelMulti,
} from '../../components/device-label';
import {
  Button,
  Row,
  Text,
  Select,
  MultiSelect,
  Checkbox,
  Icon,
  toaster,
} from '../../components/core';
import ServiceMetricsForm from './service-metrics-form';
import ServiceMetricsSettings from './service-metrics-settings';
import { labelColor } from '../../helpers/labels';

const Service = ({
  route: {
    data: { params, applications, metrics, devices },
  },
}) => {
  const [selectValue, setSelectValue] = useState(() => {
    const storedValue = storage.get('selectedService', params.project);
    if (storedValue) {
      return JSON.stringify(storedValue);
    } else {
      return null;
    }
  });
  const [showMetricsForm, setShowMetricsForm] = useState();
  const [showSettings, setShowSettings] = useState();
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
            color: labelColor(label),
          },
        }),
        []
      ),
    [devices]
  );

  const navigation = useNavigation();

  const hideMetricsForm = () => setShowMetricsForm(false);
  const hideSettings = () => setShowSettings(false);
  const clearMetricToDelete = () => setMetricToDelete(null);

  const selection = selectValue && JSON.parse(selectValue);

  let selectedMetrics = [];

  if (selection && selection.application && selection.service) {
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
          if (
            m.applicationId === selection.application.id &&
            m.service === selection.service
          ) {
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

  const submitDelete = async () => {
    clearMetricToDelete();
    try {
      await api.updateServiceMetricsConfig({
        projectId: params.project,
        data: metrics.map(m => {
          if (
            m.applicationId === selection.application.id &&
            m.service === selection.service
          ) {
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
      toaster.success('Metric deleted.');
      navigation.refresh();
    } catch (error) {
      toaster.danger('Metric deletion failed.');
      console.error(error);
    }
  };

  useEffect(() => {
    storage.set('selectedService', selection, params.project);
  }, [selection]);

  const tableData = useMemo(
    () => selectedMetrics.filter(({ name }) => !!name),
    [selectedMetrics]
  );

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
            <MultiSelect
              multi
              options={labelsOptions}
              value={editRow.labels.map(label => ({
                label,
                value: label,
                props: { color: labelColor(label) },
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
                  color={labelColor(label)}
                />
              ))}
            </Row>
          ),
      },
      {
        id: 'device',
        accessor: ({ properties }) =>
          properties && properties.includes('device'),
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
          value: JSON.stringify({ application, service }),
        })),
    [applications]
  );

  let metricEndpointConfigs;
  if (selection && selection.application) {
    const app = applications.find(({ id }) => id === selection.application.id);
    if (app) {
      metricEndpointConfigs = app.metricEndpointConfigs;
    }
  }

  return (
    <>
      <Row marginBottom={4} width={11}>
        <Select
          onChange={e => setSelectValue(e.target.value)}
          value={selectValue}
          options={selectOptions}
          placeholder="Select a Service"
          none="No services"
        />
      </Row>
      <Card
        title="Service Metrics"
        subtitle="These are custom metrics you define on your services."
        size="full"
        actions={[
          {
            title: 'Settings',
            variant: 'secondary',
            onClick: () => setShowSettings(true),
          },
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
      <Popup show={!!showSettings} onClose={hideSettings} overflow="visible">
        <ServiceMetricsSettings
          projectId={params.project}
          applicationId={
            selection && selection.application && selection.application.id
          }
          service={selection && selection.service}
          metricEndpointConfigs={metricEndpointConfigs}
          close={hideSettings}
        />
      </Popup>
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
        />
      </Popup>
      <Popup show={!!metricToDelete} onClose={clearMetricToDelete}>
        <Card border title="Delete Service Metric" size="large">
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
