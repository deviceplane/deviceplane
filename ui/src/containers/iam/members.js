import React, { useMemo } from 'react';

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
      { Header: 'Email', accessor: 'user.email' },
      {
        Header: 'Name',
        accessor: ({ user: { firstName, lastName } }) =>
          `${firstName} ${lastName}`,
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

  return (
    <Card
      title="Members"
      size="xlarge"
      actions={[{ href: 'add', title: 'Add member' }]}
      maxHeight="100%"
    >
      <Table
        columns={columns}
        data={tableData}
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
