import React, { useMemo } from 'react';
import moment from 'moment';

import { renderLabels } from '../helpers/labels';
import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Text } from '../components/core';

const Provisioning = ({
  route: {
    data: { params, registrationTokens },
  },
}) => {
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name', minWidth: '150px' },
      {
        Header: 'Created At',
        accessor: ({ createdAt }) =>
          createdAt ? moment(createdAt).fromNow() : '-',
        minWidth: '150px',
      },
      {
        Header: 'Registered Devices',
        accessor: 'deviceCounts.allCount',
        minWidth: '130px',
        maxWidth: '130px',
      },
      {
        Header: 'Registration Limit',
        accessor: 'maxRegistrations',
        Cell: ({ row: { original } }) => (
          <Text>
            {typeof original.maxRegistrations === 'number'
              ? original.maxRegistrations
              : 'Unlimited'}
          </Text>
        ),
        minWidth: '130px',
        maxWidth: '130px',
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ row: { original } }) =>
          original.labels ? renderLabels(original.labels) : null,
        minWidth: '200px',
        maxWidth: '2fr',
        cellStyle: {
          marginBottom: '-8px',
        },
      },
    ],
    []
  );
  const tableData = useMemo(() => registrationTokens, [registrationTokens]);

  return (
    <Layout>
      <Card
        title="Registration Tokens"
        size="full"
        maxHeight="100%"
        actions={[
          {
            href: 'registration-tokens/create',
            title: 'Create Registration Token',
          },
        ]}
      >
        <Table
          data={tableData}
          columns={columns}
          rowHref={({ name }) =>
            `/${params.project}/provisioning/registration-tokens/${name}`
          }
          placeholder={
            <Text>
              There are no <strong>Registration Tokens</strong>.
            </Text>
          }
        />
      </Card>
    </Layout>
  );
};

export default Provisioning;
