import React from 'react';

import { Row, Text } from '../components/core';

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
      {Object.keys(labels).map((key, i) => (
        <Row
          marginY={2}
          marginRight={2}
          overflow="hidden"
          key={key}
          fontSize={0}
          fontWeight={3}
          onClick={e => {
            e.stopPropagation();
            onClick(key, labels[key]);
          }}
        >
          <Text
            backgroundColor={labelColorMap[key]}
            paddingX={2}
            paddingY={1}
            color="black"
            borderTopLeftRadius={1}
            borderBottomLeftRadius={1}
            whiteSpace="nowrap"
            overflow="hidden"
          >
            {key}
          </Text>
          <Text
            backgroundColor="grays.3"
            paddingX={2}
            paddingY={1}
            borderTopRightRadius={1}
            borderBottomRightRadius={1}
            overflow="hidden"
            whiteSpace="nowrap"
          >
            {labels[key]}
          </Text>
        </Row>
      ))}
    </Row>
  );
}
