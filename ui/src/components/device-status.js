import React from 'react';
import moment from 'moment';

import { Row, Column, Text } from './core';

export const StatusOffline = 'offline';
export const StatusOnline = 'online';

const getColor = status => {
  switch (status) {
    case StatusOnline:
      return 'green';
    case StatusOffline:
    default:
      return 'red';
  }
};

const DeviceStatus = ({ status, lastSeenAt, inline }) => {
  const statusText = (
    <Text
      style={{ textTransform: 'uppercase' }}
      color={getColor(status)}
      fontSize={0}
      fontWeight={3}
    >
      {status}
    </Text>
  );

  if (status === StatusOffline && lastSeenAt) {
    if (inline) {
      return (
        <Row alignItems="center">
          {statusText}
          <Text fontSize={0} fontWeight={2} color="grays.8" marginLeft={2}>
            ({' '}
            <Text fontSize={0} fontWeight={1} color="grays.12">
              {moment(lastSeenAt).toNow(true)}
            </Text>{' '}
            )
          </Text>
        </Row>
      );
    }
    return (
      <Column>
        {statusText}
        <Text fontSize={0} color="grays.11">
          {moment(lastSeenAt).toNow(true)}
        </Text>
      </Column>
    );
  }

  return statusText;
};

export default DeviceStatus;
