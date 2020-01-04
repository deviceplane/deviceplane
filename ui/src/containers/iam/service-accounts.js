import React, { useMemo } from 'react';
import { useNavigation } from 'react-navi';

import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const ServiceAccounts = ({
  route: {
    data: { params, serviceAccounts },
  },
}) => {
  const navigation = useNavigation();
  const columns = useMemo(
    () => [
      { Header: 'Name', accessor: 'name' },
      { Header: 'Description', accessor: 'description' },
      {
        Header: 'Roles',
        Cell: ({
          row: {
            original: { roles },
          },
        }) => <Text>{roles.map(({ name }) => name).join(', ')}</Text>,
      },
    ],
    []
  );
  const tableData = useMemo(() => serviceAccounts, [serviceAccounts]);

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
        columns={columns}
        data={tableData}
        onRowSelect={({ name }) =>
          navigation.navigate(`/${params.project}/iam/service-accounts/${name}`)
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
