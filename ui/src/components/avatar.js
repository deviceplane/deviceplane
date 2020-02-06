import React from 'react';

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

const Avatar = ({ name, color = 'white', borderColor = 'white' }) => (
  <Column
    width="32px"
    height="32px"
    alignItems="center"
    justifyContent="center"
    bg="black"
    borderRadius={6}
    border={0}
    borderColor={borderColor}
  >
    <Text color={color} fontWeight={2} fontSize={0}>
      {getInitials(name)}
    </Text>
  </Column>
);

export default Avatar;
