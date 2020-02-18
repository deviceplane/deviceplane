import React from 'react';

import { Text } from './core';

export const ServiceStateUnknown = 'unknown';
export const ServiceStatePullingImage = 'pulling image';
export const ServiceStateCreatingContainer = 'creating container';
export const ServiceStateStoppingPreviousContainer =
  'stopping previous container';
export const ServiceStateRemovingPreviousContainer =
  'removing previous container';
export const ServiceStateStartingContainer = 'starting container';
export const ServiceStateRunning = 'running';
export const ServiceStateExited = 'exited';

const getColor = state => {
  switch (state) {
    case ServiceStateRunning:
      return 'green';
    case ServiceStateExited:
      return 'red';
    default:
      return 'white';
  }
};

const ServiceState = ({ state }) => {
  return (
    <Text
      style={{ textTransform: 'uppercase' }}
      fontSize={0}
      fontWeight={3}
      color={getColor(state)}
    >
      {state}
    </Text>
  );
};

export default ServiceState;
