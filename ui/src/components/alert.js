import React from 'react';
import { Icon } from 'evergreen-ui';

import theme from '../theme';
import { Column, Row, Text } from './core';

const getIcon = variant => {
  switch (variant) {
    case 'success':
      return 'endorsed';
    case 'error':
      return 'error';
    case 'info':
    default:
      return 'info-sign';
  }
};

const getColor = variant => {
  switch (variant) {
    case 'success':
      return theme.colors.green;
    case 'error':
      return theme.colors.red;
    case 'info':
    default:
      return theme.colors.primary;
  }
};

const Alert = ({ show, title, description, variant = 'info', children }) => {
  if (!show) {
    return null;
  }
  const color = getColor(variant);
  return (
    <Column
      bg="black"
      border={0}
      borderColor={color}
      borderRadius={1}
      padding={4}
      marginBottom={5}
    >
      <Column marginBottom={children ? 4 : 0}>
        <Row alignItems="center">
          <Icon
            icon={getIcon(variant)}
            color={color}
            size={14}
            flexShrink={0}
          />
          {title ? (
            <Text fontSize={4} fontWeight={2} marginLeft={2}>
              {title}
            </Text>
          ) : (
            <Text fontWeight={1} marginLeft={2}>
              {description}
            </Text>
          )}
        </Row>

        {title && description && (
          <Text color="grays.8" marginTop={3} fontWeight={1}>
            {description}
          </Text>
        )}
      </Column>

      {children}
    </Column>
  );
};

export default Alert;
