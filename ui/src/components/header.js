import React from 'react';

import { Row } from './core';
import AvatarMenu from './avatar-menu';

const Header = ({ children }) => {
  return (
    <Row
      alignItems="center"
      justifyContent="space-between"
      alignSelf="stretch"
      padding={5}
      paddingBottom={1}
    >
      <Row flex={1} alignItems="center">
        {/* <Text fontSize={3} fontWeight={4}>
          {route.data.params.project}
        </Text> */}
      </Row>
      <Row justifyContent="center" flex={1}>
        {children}
      </Row>
      <Row justifyContent="flex-end" alignItems="flex-start" flex={1}>
        <AvatarMenu />
      </Row>
    </Row>
  );
};

export default Header;
