import React, { useMemo } from 'react';

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
      },
      {
        Header: 'Applications',
        accessor: 'applicationCounts.allCount',
      },
    ],
    []
  );
  const tableData = useMemo(() => projects, [projects]);

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
          columns={columns}
          data={tableData}
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
