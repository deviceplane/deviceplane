import React, { useMemo } from 'react';

import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const Roles = ({
  route: {
    data: { params, roles },
  },
}) => {
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name' },
      { Header: 'Description', accessor: 'description', maxWidth: '2fr' },
    ],
    []
  );
  const tableData = useMemo(() => roles, [roles]);

  return (
    <Card
      title="Roles"
      size="xlarge"
      maxHeight="100%"
      actions={[
        {
          href: `create`,
          title: 'Create role',
        },
      ]}
    >
      <Table
        columns={columns}
        data={tableData}
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
