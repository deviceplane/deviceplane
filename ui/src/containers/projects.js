import React, { useMemo } from 'react';
import { useNavigation } from 'react-navi';

import Layout from '../components/layout';
import Card from '../components/card';
import Table from '../components/table';
import { Text } from '../components/core';

const Projects = ({
  route: {
    data: { projects },
  },
}) => {
  const navigation = useNavigation();
  const columns = useMemo(
    () => [
      {
        Header: 'Name',
        accessor: 'name',
        style: {
          flex: 2,
        },
      },
      {
        Header: 'Devices',
        accessor: 'deviceCounts.allCount',
        style: {
          flex: '0 0 76px',
        },
      },
      {
        Header: 'Applications',
        accessor: 'applicationCounts.allCount',
        style: {
          flex: '0 0 120px',
        },
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
