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
    title: 'Releases',
    to: 'releases',
  },
  {
    title: 'Scheduling',
    to: 'scheduling',
  },
  {
    title: 'Release Pinning',
    to: 'release-pinning',
  },
  {
    title: 'Settings',
    to: 'settings',
  },
];

const Application = ({ route }) => {
  if (!route) {
    return null;
  }

  return (
    <Layout
      header={
        <Tabs
          content={tabs.map(({ to, title }) => ({
            title,
            href: `/${route.data.params.project}/applications/${route.data.application.name}/${to}`,
          }))}
        />
      }
      alignItems="center"
    >
      <View />
    </Layout>
  );
};

export default Application;
