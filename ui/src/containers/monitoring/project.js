import React, { useMemo } from 'react';
import { useTable, useSortBy, useRowSelect } from 'react-table';

import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Table, {
  SelectColumn,
  DeviceLabelKeyColumn,
} from '../../components/table';
import Popup from '../../components/popup';
import { Column, Row, Text, Checkbox, Button } from '../../components/core';
import ProjectMetricOverview from './project-metric-overview';

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
  const [isMetricOverview, toggleMetricOverview] = useToggle();

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

  const columns = useMemo(
    () => [
      SelectColumn,
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
      DeviceLabelKeyColumn,
      {
        id: 'device',
        accessor: ({ properties }) =>
          properties && properties.includes('device'),
        Header: (
          <Row title="When enabled, a Datadog tag with the device name is included.">
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
      {
        Header: 'Enabled',
        accessor: 'enabled',
        Cell: ({ cell: { value } }) => <Checkbox readOnly checked={value} />,
        minWidth: '110px',
        maxWidth: '110px',
        cellStyle: {
          justifyContent: 'center',
        },
      },
    ],
    []
  );

  const { selectedFlatRows, ...tableProps } = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy,
    useRowSelect
  );

  return (
    <>
      <Card
        title="Project Metrics"
        size="full"
        subtitle="These metrics provide insights over your entire project."
      >
        <Row marginBottom={3}>
          <Button
            marginRight={4}
            title="Edit"
            variant="tertiary"
            disabled={selectedFlatRows.length !== 1}
            onClick={toggleMetricOverview}
          />
        </Row>
        <Table {...tableProps} />
      </Card>
      <Popup show={isMetricOverview} onClose={toggleMetricOverview}>
        <ProjectMetricOverview
          projectId={params.project}
          devices={devices}
          metrics={supportedMetrics}
          metric={selectedFlatRows[0] && selectedFlatRows[0].original}
          close={toggleMetricOverview}
        />
      </Popup>
    </>
  );
};

export default Project;
