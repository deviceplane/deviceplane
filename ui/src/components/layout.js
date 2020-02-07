import React from 'react';

import { Column, Row } from './core';
import Header from './header';
import Sidebar from './sidebar';

const Layout = ({ children, header, ...rest }) => (
  <Column height="100%">
    <Header>{header}</Header>
    <Row height="100%">
      <Sidebar />
      <Column flex={1} overflow="auto">
        <Column
          flex={1}
          {...rest}
          padding={5}
          overflow="auto"
          bg="pageBackground"
          borderRadius={1}
        >
          {children}
        </Column>
      </Column>
    </Row>
  </Column>
);

export default Layout;
