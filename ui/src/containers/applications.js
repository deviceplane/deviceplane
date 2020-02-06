import React, { useMemo } from 'react';
import moment from 'moment';

import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Column, Text } from '../components/core';

const Applications = ({
  route: {
    data: { params, applications },
  },
}) => {
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name', style: { flex: 2 } },
      {
        Header: 'Services',
        accessor: 'latestRelease.config',
        Cell: ({ cell: { value: config } }) =>
          config ? (
            <Column>
              {Object.keys(config).map(name => (
                <Text>{name}</Text>
              ))}
            </Column>
          ) : (
            '-'
          ),
        style: { flex: 2 },
      },
      {
        Header: 'Scheduling',
        accessor: ({ schedulingRule: { scheduleType } }) =>
          scheduleType === 'AllDevices'
            ? 'All devices'
            : scheduleType === 'NoDevices'
            ? 'No devices'
            : 'Conditional',
        style: { flex: '0 0 120px' },
      },
      {
        Header: 'Last Release',
        accessor: 'latestRelease.createdAt',
        Cell: ({ cell: { value } }) => (
          <Text>{value ? moment(value).fromNow() : '-'}</Text>
        ),
        style: {
          flex: '0 0 150px',
        },
      },
      {
        Header: 'Devices',
        accessor: 'deviceCounts.allCount',
        style: {
          flex: '0 0 76px',
        },
      },
    ],
    []
  );
  const tableData = useMemo(() => applications, [applications]);
  return (
    <Layout alignItems="center">
      <Card
        title="Applications"
        size="xxlarge"
        actions={[{ title: 'Create Application', href: 'create' }]}
        maxHeight="100%"
      >
        <Table
          columns={columns}
          data={tableData}
          rowHref={({ name }) => `/${params.project}/applications/${name}`}
          placeholder={
            <Text>
              There are no <strong>Applications</strong>.
            </Text>
          }
        />
      </Card>
    </Layout>
  );
};

export default Applications;
