import React from 'react';

import theme from '../theme';
import { Row, Text } from './core';

const DeviceLabel = ({
  label: { key, value },
  color = theme.colors.white,
  onClick = () => {},
}) => {
  return (
    <Row
      flex={1}
      fontSize={0}
      fontWeight={3}
      onClick={e => {
        e.stopPropagation();
        onClick({ key, value });
      }}
    >
      <Text
        backgroundColor={color}
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
        {value}
      </Text>
    </Row>
  );
};

export const DeviceLabelMulti = ({ children, color }) => (
  <Row
    overflow="hidden"
    bg={color}
    borderTopLeftRadius={1}
    borderBottomLeftRadius={1}
    paddingX={2}
    paddingY={1}
  >
    <Text
      color="black"
      whiteSpace="nowrap"
      overflow="hidden"
      fontSize={0}
      fontWeight={3}
    >
      {children}
    </Text>
  </Row>
);

export const DeviceLabelKey = ({ label, color }) => (
  <Row
    marginY={2}
    marginRight={2}
    overflow="hidden"
    bg={color}
    borderRadius={1}
    paddingX={2}
    paddingY={1}
  >
    <Text
      color="black"
      whiteSpace="nowrap"
      overflow="hidden"
      fontSize={0}
      fontWeight={3}
    >
      {label}
    </Text>
  </Row>
);

export default DeviceLabel;
