import React from 'react';

import { Column, Text, Heading } from './core';

const NotFound = () => {
  return (
    <Column
      alignItems="center"
      justifyContent="center"
      flex={1}
      padding={8}
      marginTop={9}
    >
      <Heading variant="primary" fontSize={8} fontWeight={4}>
        404
      </Heading>
      <Text marginTop={4} fontSize={6} color="white">
        We couldn’t find the page you were looking for.
      </Text>
    </Column>
  );
};

export default NotFound;
