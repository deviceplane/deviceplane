import React from 'react';

import { labelColor } from '../helpers/labels';
import { Row, Text } from './core';

const DeviceLabel = ({ label: { key, value }, onClick = () => {} }) => {
  const color = labelColor(key);
  return (
    <Row
      overflow="hidden"
      marginRight={2}
      marginY={2}
      onClick={e => {
        e.stopPropagation();
        onClick({ key, value });
      }}
      border={0}
      borderColor={color}
      borderRadius={1}
    >
      <Text
        paddingX={2}
        paddingY={1}
        color={color}
        whiteSpace="nowrap"
        overflow="hidden"
        fontSize={0}
        fontWeight={2}
        borderRight={0}
        borderColor={color}
      >
        {key}
      </Text>
      <Text
        paddingX={2}
        paddingY={1}
        overflow="hidden"
        whiteSpace="nowrap"
        fontSize={0}
        fontWeight={1}
      >
        {value}
      </Text>
    </Row>
  );
};

export const DeviceLabelMulti = ({ children, color }) => (
  <Row
    overflow="hidden"
    border={0}
    borderColor={color}
    borderTopLeftRadius={1}
    borderBottomLeftRadius={1}
    paddingX={2}
    paddingY={1}
  >
    <Text
      color={color}
      whiteSpace="nowrap"
      overflow="hidden"
      fontSize={0}
      fontWeight={2}
    >
      {children}
    </Text>
  </Row>
);

export const DeviceLabelKey = ({ label }) => {
  const color = labelColor(label);
  return (
    <Row
      marginY={2}
      marginRight={2}
      overflow="hidden"
      border={0}
      borderColor={color}
      borderRadius={1}
      paddingX={2}
      paddingY={1}
    >
      <Text
        color={color}
        whiteSpace="nowrap"
        overflow="hidden"
        fontSize={0}
        fontWeight={2}
      >
        {label}
      </Text>
    </Row>
  );
};

export default DeviceLabel;
