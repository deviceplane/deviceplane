import React from 'react';

import { labelColors } from '../theme';
import DeviceLabel from '../components/device-label';
import { Row } from '../components/core';

export const renderLabels = (labels, onClick = () => {}) => (
  <Row flexWrap="wrap" overflow="hidden">
    {Object.keys(labels).map(key => (
      <DeviceLabel
        key={key}
        label={{ key, value: labels[key] }}
        color={labelColor(key)}
        onClick={onClick}
      />
    ))}
  </Row>
);

const strHash = str => {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash += str.charCodeAt(i);
  }
  return hash;
};

export const labelColor = name => {
  const hash = strHash(name);
  console.log(hash);
  const index = hash % labelColors.length;
  return labelColors[index];
};
