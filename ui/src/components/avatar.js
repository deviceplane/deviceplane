import React from 'react';
import styled from 'styled-components';

import { Column, Text } from './core';

const getInitials = (name, fallback = '?') => {
  if (!name || typeof name !== 'string') return fallback;
  return name
    .replace(/\s+/, ' ')
    .split(' ') // Repeated spaces results in empty strings
    .slice(0, 2)
    .map(v => v && v[0].toUpperCase()) // Watch out for empty strings
    .join('');
};

const Container = styled(Column)`
  &:hover {
    border-color: ${props => props.theme.colors.primary};
  }
  &:hover span {
    color: ${props => props.theme.colors.primary};
  }
`;

const Avatar = ({ name, color = 'white', borderColor = 'white' }) => (
  <Container
    width="32px"
    height="32px"
    alignItems="center"
    justifyContent="center"
    bg="black"
    borderRadius="50%"
    border={0}
    borderColor={borderColor}
  >
    <Text color={color} fontWeight={2} fontSize={0}>
      {getInitials(name)}
    </Text>
  </Container>
);

export default Avatar;
