import React, { useMemo } from 'react';
import moment from 'moment';

import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const Releases = ({
  route: {
    data: { params, application, releases },
  },
}) => {
  const columns = useMemo(
    () => [
      {
        Header: 'Release',
        accessor: 'number',
      },
      {
        Header: 'Released by',
        accessor: ({ createdByUser, createdByServiceAccount }) => {
          if (createdByUser) {
            return `${createdByUser.firstName} ${createdByUser.lastName}`;
          } else if (createdByServiceAccount) {
            return createdByServiceAccount.name;
          }
          return '-';
        },
      },
      {
        Header: 'Started At',
        accessor: ({ createdAt }) =>
          createdAt ? moment(createdAt).fromNow() : '-',
      },
      {
        Header: 'Devices',
        accessor: 'deviceCounts.allCount',
      },
    ],
    []
  );
  const tableData = useMemo(() => releases, [releases]);

  return (
    <Card
      title="Releases"
      size="xlarge"
      actions={[
        {
          title: 'Create Release',
          href: `/${params.project}/applications/${application.name}/releases/create`,
        },
      ]}
      maxHeight="100%"
    >
      <Table
        columns={columns}
        data={tableData}
        rowHref={({ id }) =>
          `/${params.project}/applications/${application.name}/releases/${id}`
        }
        placeholder={
          <Text>
            There are no <strong>Releases</strong>.
          </Text>
        }
      />
    </Card>
  );
};

export default Releases;
