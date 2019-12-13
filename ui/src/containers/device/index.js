import React from 'react';
import { View } from 'react-navi';

import Layout from '../../components/layout';
import Tabs from '../../components/tabs';

const tabs = [
  {
    title: 'Overview',
    to: 'overview',
  },
  {
    title: 'SSH',
    to: 'ssh',
  },
  {
    title: 'Settings',
    to: 'settings',
  },
];

const Device = ({ route }) => {
  if (!route) {
    return null;
  }
  return (
    <Layout
      alignItems="center"
      header={
        <Tabs
          content={tabs.map(({ to, title }) => ({
            title,
            href: `/${route.data.params.project}/devices/${route.data.device.name}/${to}`,
          }))}
        />
      }
    >
      <View />
    </Layout>
  );
};

export default Device;
