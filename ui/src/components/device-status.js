import React from 'react';

import { Badge } from './core';

const DeviceStatus = ({ status }) => {
  switch (status) {
    case 'online':
      return <Badge bg="green">online</Badge>;
    case 'offline':
    default:
      return <Badge bg="red">offline</Badge>;
  }
};

export default DeviceStatus;
