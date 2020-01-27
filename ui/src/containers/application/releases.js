import React, { useMemo } from 'react';
import moment from 'moment';
import { useNavigation } from 'react-navi';

import Card from '../../components/card';
import Table from '../../components/table';
import { Text } from '../../components/core';

const Releases = ({
  route: {
    data: { params, application, releases },
  },
}) => {
  const navigation = useNavigation();
  const columns = useMemo(
    () => [
      {
        Header: 'Release',
        accessor: 'number',
        style: {
          flex: '0 0 60px',
        },
      },
      {
        Header: 'Released by',
        accessor: ({ release }) => {
          if (release) {
            if (release.createdByUser) {
              return `${release.createdByUser.firstName} ${release.createdByUser.lastName}`;
            } else if (release.createdByServiceAccount) {
              return release.createdByServiceAccount.name;
            }
          }
          return '-';
        },
      },
      {
        Header: 'Started',
        accessor: 'createdAt',
        Cell: ({
          row: {
            original: { createdAt },
          },
        }) => <Text>{moment(createdAt).fromNow()}</Text>,
        style: {
          flex: '0 0 150px',
        },
      },
      {
        Header: 'Device count',
        accessor: 'deviceCounts.allCount',
        style: {
          flex: '0 0 102px',
        },
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
        onRowSelect={({ id }) =>
          navigation.navigate(
            `/${params.project}/applications/${application.name}/releases/${id}`
          )
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
