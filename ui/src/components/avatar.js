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
    width="36px"
    height="36px"
    alignItems="center"
    justifyContent="center"
    bg="black"
    borderRadius={6}
    border={0}
    borderColor="white"
  >
    <Text color="white" fontWeight={3} fontSize={0}>
      {getInitials(name)}
    </Text>
  </Column>
);

export default Avatar;
