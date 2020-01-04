import React, { useMemo } from 'react';
import moment from 'moment';
import { useNavigation } from 'react-navi';

import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Text } from '../components/core';

const Applications = ({
  route: {
    data: { params, applications },
  },
}) => {
  const navigation = useNavigation();
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name', style: { flex: 2 } },
      {
        Header: 'Last Release',
        accessor: 'latestRelease.createdAt',
        Cell: ({ row: { original } }) => (
          <Text>
            {original.latestRelease
              ? moment(original.latestRelease.createdAt).fromNow()
              : '-'}
          </Text>
        ),
      },
      {
        Header: 'Device Count',
        accessor: 'deviceCounts.allCount',
      },
    ],
    []
  );
  const tableData = useMemo(() => applications, [applications]);
  return (
    <Layout title="Applications" alignItems="center">
      <Card
        title="Applications"
        size="xlarge"
        actions={[{ title: 'Create Application', href: 'create' }]}
        maxHeight="100%"
      >
        <Table
          columns={columns}
          data={tableData}
          onRowSelect={row =>
            navigation.navigate(`/${params.project}/applications/${row.name}`)
          }
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
