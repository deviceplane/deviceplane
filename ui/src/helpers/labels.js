import React from 'react';

import DeviceLabel from '../components/device-label';
import { Row } from '../components/core';

export function buildLabelColorMap(oldLabelColorMap, labelColors, items) {
  var x = [];
  items.forEach(item => {
    Object.keys(item.labels).forEach(label => {
      x.push(label);
    });
  });
  const labelKeys = [...new Set(x)].sort();

  var labelColorMap = Object.assign({}, oldLabelColorMap);
  labelKeys.forEach((key, i) => {
    if (!labelColorMap[key]) {
      labelColorMap[key] = labelColors[i % (labelColors.length - 1)];
    }
  });
  return labelColorMap;
}

export function renderLabels(labels, labelColorMap, onClick = () => {}) {
  return (
    <Row flexWrap="wrap" overflow="hidden">
      {Object.keys(labels).map(key => (
        <DeviceLabel
          key={key}
          label={{ key, value: labels[key] }}
          color={labelColorMap[key]}
          onClick={onClick}
        />
      ))}
    </Row>
  );
}
