import React, { useMemo } from 'react';
import { useTable, useSortBy } from 'react-table';

import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Text } from '../components/core';

const Projects = ({
  route: {
    data: { projects },
  },
}) => {
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
    ],
    []
  );
  const tableData = useMemo(() => projects, [projects]);

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
        maxHeight="100%"
        actions={[
          {
            href: '/projects/create',
            title: 'Create Project',
          },
        ]}
      >
        <Table
          {...tableProps}
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
