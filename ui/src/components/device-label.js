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
      overflow="hidden"
      marginRight={2}
      marginBottom={2}
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
        fontSize={0}
        fontWeight={2}
      >
        {key}
      </Text>
      <Text
        backgroundColor="grays.2"
        paddingX={2}
        paddingY={1}
        borderTopRightRadius={1}
        borderBottomRightRadius={1}
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
      fontWeight={2}
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
      fontWeight={2}
    >
      {label}
    </Text>
  </Row>
);

export default DeviceLabel;
