import React, { useMemo, useState, useEffect } from 'react';
import { useNavigation } from 'react-navi';
import { useTable, useSortBy, useRowSelect } from 'react-table';

import storage from '../../storage';
import api from '../../api';
import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Table, {
  SelectColumn,
  DeviceLabelKeyColumn,
} from '../../components/table';
import Popup from '../../components/popup';
import {
  Button,
  Row,
  Text,
  Select,
  Checkbox,
  toaster,
} from '../../components/core';
import ServiceMetricsForm from './service-metrics-form';
import ServiceMetricsSettings from './service-metrics-settings';
import MetricOverview from './metric-overview';

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
  const [isMetricsForm, toggleMetricsForm] = useToggle();
  const [isSettings, toggleSettings] = useToggle();
  const [isDeleteForm, toggleDeleteForm] = useToggle();
  const [isMetricOverview, toggleMetricOverview] = useToggle();
  const navigation = useNavigation();

  const selection = useMemo(() => selectValue && JSON.parse(selectValue), [
    selectValue,
  ]);

  const submitDelete = async () => {
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
                ({ name }) =>
                  !selectedFlatRows.find(
                    ({ original }) => original.name === name
                  )
              ),
            };
          }
          return m;
        }),
      });
      toggleDeleteForm();
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

  const columns = useMemo(
    () => [
      SelectColumn,
      {
        Header: 'Metric',
        accessor: 'name',
      },
      DeviceLabelKeyColumn,
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
        Cell: ({ cell: { value } }) => <Checkbox readOnly checked={value} />,
        minWidth: '100px',
        maxWidth: '100px',
        cellStyle: {
          justifyContent: 'center',
        },
      },
    ],
    []
  );

  const selectedMetrics = useMemo(() => {
    if (selection && selection.application && selection.service) {
      const data = metrics.find(
        ({ applicationId, service }) =>
          applicationId === selection.application.id &&
          service === selection.service
      );
      return data ? data.exposedMetrics.filter(({ name }) => !!name) : [];
    }
    return [];
  }, [metrics, selection]);

  const { selectedFlatRows, ...tableProps } = useTable(
    {
      columns,
      data: selectedMetrics,
    },
    useSortBy,
    useRowSelect
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
            onClick: toggleSettings,
          },
          {
            title: 'Add Service Metrics',
            onClick: toggleMetricsForm,
          },
        ]}
        disabled={!(selection && selection.service)}
        maxHeight="100%"
      >
        <Row marginBottom={3}>
          <Button
            marginRight={4}
            title="Edit"
            variant="tertiary"
            disabled={selectedFlatRows.length !== 1}
            onClick={toggleMetricOverview}
          />
          <Button
            title="Delete"
            variant="tertiaryDanger"
            disabled={selectedFlatRows.length === 0}
            onClick={toggleDeleteForm}
          />
        </Row>
        <Table
          {...tableProps}
          placeholder={
            <Text>
              There are no <strong>Service Metrics</strong>.
            </Text>
          }
        />
      </Card>
      <Popup show={isMetricOverview} onClose={toggleMetricOverview}>
        <MetricOverview
          service={selection && selection.service}
          application={selection && selection.application}
          projectId={params.project}
          devices={devices}
          metrics={selectedMetrics}
          metric={selectedFlatRows[0] && selectedFlatRows[0].original}
          close={toggleMetricOverview}
        />
      </Popup>
      <Popup show={isSettings} onClose={toggleSettings} overflow="visible">
        <ServiceMetricsSettings
          projectId={params.project}
          applicationId={
            selection && selection.application && selection.application.id
          }
          service={selection && selection.service}
          metricEndpointConfigs={metricEndpointConfigs}
          close={toggleSettings}
        />
      </Popup>
      <Popup
        show={isMetricsForm}
        onClose={toggleMetricsForm}
        overflow="visible"
      >
        <ServiceMetricsForm
          params={params}
          allMetrics={metrics}
          metrics={selectedMetrics}
          devices={devices}
          application={selection && selection.application}
          service={selection && selection.service}
          close={toggleMetricsForm}
        />
      </Popup>
      <Popup show={isDeleteForm} onClose={toggleDeleteForm}>
        <Card
          border
          title={`Delete Service Metric${
            selectedFlatRows.length > 1 ? 's' : ''
          }`}
          size="large"
        >
          <Text>
            You are about to delete the{' '}
            <strong>
              {selectedFlatRows
                .map(({ original: { name } }) => name)
                .join(', ')}
            </strong>{' '}
            metric{selectedFlatRows.length > 1 ? 's' : ''}.
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
