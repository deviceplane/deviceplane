import React, { useMemo } from 'react';
import { useTable, useSortBy } from 'react-table';

import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const Members = ({
  route: {
    data: { params, members },
  },
}) => {
  const columns = useMemo(
    () => [
      {
        Header: 'Email',
        accessor: ({ user: { email, providerName } }) =>
          `${email}${providerName ? ` (via ${providerName})` : ''}`,
      },
      {
        Header: 'Name',
        accessor: 'user.name',
      },
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
  const tableData = useMemo(() => members, [members]);

  const tableProps = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy
  );

  return (
    <Card
      title="Members"
      size="xlarge"
      actions={[{ href: 'add', title: 'Add member' }]}
      maxHeight="100%"
    >
      <Table
        {...tableProps}
        rowHref={({ user: { id } }) => `/${params.project}/iam/members/${id}`}
        placeholder={
          <Text>
            There are no <strong>Members</strong>.
          </Text>
        }
      />
    </Card>
  );
};

export default Members;
