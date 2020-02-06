import React from 'react';
import styled from 'styled-components';

import AvatarMenu from './avatar-menu';
import ProjectSelector from './project-selector';
import { Row, MenuItem } from './core';

const Header = ({ children }) => {
  return (
    <Row
      alignItems="center"
      justifyContent="space-between"
      alignSelf="stretch"
      padding={5}
      paddingTop={4}
      paddingBottom={3}
    >
      <Row flex={1} alignItems="center">
        <ProjectSelector />
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
