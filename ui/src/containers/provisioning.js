import React, { useState, useMemo } from 'react';
import moment from 'moment';
import { useNavigation } from 'react-navi';

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
  const navigation = useNavigation();
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name' },
      {
        Header: 'Created At',
        accessor: 'createdAt',
        Cell: ({ row: { original } }) => (
          <Text>
            {original.createdAt ? moment(original.createdAt).fromNow() : '-'}
          </Text>
        ),
        style: {
          flex: '0 0 150px',
        },
      },
      {
        Header: 'Devices Registered',
        accessor: 'deviceCounts.allCount',
        style: { flex: '0 0 150px' },
      },
      {
        Header: 'Registration Limit',
        accessor: 'maxRegistrations',
        style: { flex: '0 0 140px' },
        Cell: ({ row: { original } }) => (
          <Text>
            {typeof original.maxRegistrations === 'number'
              ? original.maxRegistrations
              : 'Unlimited'}
          </Text>
        ),
      },
      {
        Header: 'Labels',
        accessor: 'labels',
        Cell: ({ row: { original } }) =>
          original.labels ? renderLabels(original.labels) : null,
        style: {
          flex: 2,
        },
      },
    ],
    []
  );
  const tableData = useMemo(() => registrationTokens, [registrationTokens]);

  return (
    <Layout title="Provisioning">
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
          onRowSelect={({ name }) =>
            navigation.navigate(
              `/${params.project}/provisioning/registration-tokens/${name}`
            )
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
