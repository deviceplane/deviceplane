import React, { useMemo } from 'react';
import { useTable, useSortBy } from 'react-table';

import { useRequest, endpoints } from '../../api';
import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const Roles = ({
  route: {
    data: { params },
  },
}) => {
  const { data: roles } = useRequest(
    endpoints.roles({ projectId: params.project })
  );
  const tableData = useMemo(() => roles, [roles]);
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name' },
      { Header: 'Description', accessor: 'description', maxWidth: '2fr' },
    ],
    []
  );
  const tableProps = useTable(
    {
      columns,
      data: tableData,
    },
    useSortBy
  );

  return (
    <Card
      title="Roles"
      size="xlarge"
      actions={[
        {
          href: `create`,
          title: 'Create role',
        },
      ]}
    >
      <Table
        {...tableProps}
        rowHref={({ name }) => `/${params.project}/iam/roles/${name}`}
        placeholder={
          <Text>
            There are no <strong>Roles</strong>.
          </Text>
        }
      />
    </Card>
  );
};

export default Roles;
