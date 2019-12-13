import React from 'react';
import { View } from 'react-navi';

import Layout from '../../components/layout';
import Tabs from '../../components/tabs';

const tabs = [
  {
    title: 'Members',
    to: 'iam/members',
  },
  {
    title: 'Service Accounts',
    to: 'iam/service-accounts',
  },
  {
    title: 'Roles',
    to: 'iam/roles',
  },
];

const Iam = ({ route }) => {
  if (!route) {
    return null;
  }
  return (
    <Layout
      header={
        <Tabs
          content={tabs.map(({ to, title }) => ({
            title,
            href: `/${route.data.params.project}/${to}`,
          }))}
        />
      }
      alignItems="center"
    >
      <View />
    </Layout>
  );
};

export default Iam;
