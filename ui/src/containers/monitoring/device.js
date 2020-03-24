import React, { useMemo } from 'react';
import { useNavigation } from 'react-navi';
import { useTable, useSortBy, useRowSelect } from 'react-table';

import api from '../../api';
import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Table, {
  SelectColumn,
  DeviceLabelKeyColumn,
} from '../../components/table';
import Popup from '../../components/popup';
import {
  Column,
  Row,
  Button,
  Text,
  Checkbox,
  toaster,
} from '../../components/core';
import { getMetricLabel } from '../../helpers/metrics';
import DeviceMetricsForm from './device-metrics-form';
import MetricOverview from './metric-overview';

const Device = ({
  route: {
    data: { params, devices, metrics },
  },
}) => {
  const [isDeleteForm, toggleDeleteForm] = useToggle();
  const [isMetricsForm, toggleMetricsForm] = useToggle();
  const [isMetricOverview, toggleMetricOverview] = useToggle();
  const navigation = useNavigation();

  const submitDelete = async () => {
    try {
      await api.updateDeviceMetricsConfig({
        projectId: params.project,
        data: metrics.filter(
          ({ name }) =>
            !selectedFlatRows.find(({ original }) => original.name === name)
        ),
      });
      toggleDeleteForm();
      toaster.success('Metrics deleted.');
      navigation.refresh();
    } catch (error) {
      toaster.danger('Metric deletion failed.');
      console.error(error);
    }
  };

  const columns = useMemo(
    () => [
      SelectColumn,
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

  const data = useMemo(() => metrics, [metrics]);

  const { selectedFlatRows, ...tableProps } = useTable(
    {
      columns,
      data,
    },
    useSortBy,
    useRowSelect
  );

  return (
    <>
      <Card
        title="Device Metrics"
        subtitle="These metrics provide information on the state of each device."
        size="full"
        actions={[
          {
            title: 'Add Device Metrics',
            onClick: toggleMetricsForm,
          },
        ]}
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
              There are no <strong>Device Metrics</strong>.
            </Text>
          }
        />
        <Popup show={isMetricsForm} onClose={toggleMetricsForm}>
          <DeviceMetricsForm
            params={params}
            metrics={metrics}
            devices={devices}
            close={toggleMetricsForm}
          />
        </Popup>
        <Popup show={isMetricOverview} onClose={toggleMetricOverview}>
          <MetricOverview
            projectId={params.project}
            devices={devices}
            metrics={metrics}
            metric={selectedFlatRows[0] && selectedFlatRows[0].original}
            close={toggleMetricOverview}
          />
        </Popup>
        <Popup show={isDeleteForm} onClose={toggleDeleteForm}>
          <Card
            border
            title={`Delete Device Metric${
              selectedFlatRows.length > 1 ? 's' : ''
            }`}
            size="large"
          >
            <Text>
              You are about to delete the{' '}
              <strong>
                {selectedFlatRows
                  .map(({ original: { name } }) => getMetricLabel(name))
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
      </Card>
    </>
  );
};

export default Device;
