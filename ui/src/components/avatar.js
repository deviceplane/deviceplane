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

const Avatar = ({ name }) => (
  <Column
    width={7}
    height={7}
    alignItems="center"
    justifyContent="center"
    bg="primary"
    borderRadius={6}
  >
    <Text color="black" fontWeight={2} fontSize={1}>
      {getInitials(name)}
    </Text>
  </Column>
);

export default Avatar;
