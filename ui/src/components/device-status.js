import React from 'react';

import { Text } from './core';

const getColor = status => {
  switch (status) {
    case 'online':
      return 'green';
    case 'offline':
    default:
      return 'red';
  }
};

const DeviceStatus = ({ status }) => (
  <Text
    style={{ textTransform: 'uppercase' }}
    color={getColor(status)}
    fontSize={0}
    fontWeight={3}
  >
    {status}
  </Text>
);

export default DeviceStatus;
