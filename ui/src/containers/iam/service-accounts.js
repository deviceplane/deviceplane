import React, { useMemo } from 'react';
import { useTable, useSortBy } from 'react-table';

import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const ServiceAccounts = ({
  route: {
    data: { params, serviceAccounts },
  },
}) => {
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name' },
      { Header: 'Description', accessor: 'description', maxWidth: '2fr' },
      {
        Header: 'Roles',
        accessor: 'roles',
        Cell: ({ cell: { value } }) => (
          <Text>{value.map(({ name }) => name).join(', ')}</Text>
        ),
      },
    ],
    []
  );
  const tableData = useMemo(() => serviceAccounts, [serviceAccounts]);

  const tableProps = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy
  );

  return (
    <Card
      title="Service Accounts"
      size="xlarge"
      maxHeight="100%"
      actions={[
        {
          href: `create`,
          title: 'Create service account',
        },
      ]}
    >
      <Table
        {...tableProps}
        rowHref={({ name }) =>
          `/${params.project}/iam/service-accounts/${name}`
        }
        placeholder={
          <Text>
            There are no <strong>Service Accounts</strong>.
          </Text>
        }
      />
    </Card>
  );
};

export default ServiceAccounts;
