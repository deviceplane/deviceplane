import React, { useMemo } from 'react';
import { useTable, useSortBy } from 'react-table';

import { useRequest, endpoints } from '../api';
import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Text } from '../components/core';

const Projects = () => {
  const { data } = useRequest(endpoints.projects());

  const tableData = useMemo(
    () => (data ? data.map(({ project }) => project) : []),
    [data]
  );

  const columns = useMemo(
    () => [
      {
        Header: 'Name',
        accessor: 'name',
      },
      {
        Header: 'Devices',
        accessor: 'deviceCounts.allCount',
        maxWidth: '100px',
        minWidth: '100px',
      },
      {
        Header: 'Applications',
        accessor: 'applicationCounts.allCount',
        maxWidth: '140px',
        minWidth: '140px',
      },
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
    <Layout alignItems="center">
      <Card
        title="Projects"
        size="xlarge"
        actions={[
          {
            href: '/projects/create',
            title: 'Create Project',
          },
        ]}
      >
        <Table
          {...tableProps}
          loading={!data}
          rowHref={({ name }) => `/${name}`}
          placeholder={
            <Text>
              There are no <strong>Projects</strong>.
            </Text>
          }
        />
      </Card>
    </Layout>
  );
};

export default Projects;
