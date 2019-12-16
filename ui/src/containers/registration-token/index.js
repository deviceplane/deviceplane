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
    title: 'Settings',
    to: 'settings',
  },
];

const RegistrationToken = ({ route }) => {
  if (!route) {
    return null;
  }

  return (
    <Layout
      header={
        <Tabs
          content={tabs.map(({ to, title }) => ({
            title,
            href: `/${route.data.params.project}/provisioning/registration-tokens/${route.data.registrationToken.name}/${to}`,
          }))}
        />
      }
      alignItems="center"
    >
      <View />
    </Layout>
  );
};

export default RegistrationToken;
