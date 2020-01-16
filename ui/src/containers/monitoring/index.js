import React from 'react';
import { View } from 'react-navi';

import Layout from '../../components/layout';
import Tabs from '../../components/tabs';

const tabs = [
  {
    title: 'Project',
    to: 'project',
  },
  {
    title: 'Device',
    to: 'device',
  },
  {
    title: 'Service',
    to: 'service',
  },
  {
    title: 'Integrations',
    to: 'integrations',
  },
];

const Monitoring = ({ route }) => {
  if (!route) {
    return null;
  }

  return (
    <Layout
      header={
        route.data.project.datadogApiKey ? (
          <Tabs
            content={tabs.map(({ to, title }) => ({
              title,
              href: `/${route.data.params.project}/monitoring/${to}`,
            }))}
          />
        ) : null
      }
      alignItems="center"
    >
      <View />
    </Layout>
  );
};

export default Monitoring;
