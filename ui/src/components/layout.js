import React from 'react';

import { Column, Row } from './core';
import Header from './header';
import Sidebar from './sidebar';

const Layout = ({ children, header, ...rest }) => (
  <Column height="100%">
    <Header>{header}</Header>
    <Row height="100%">
      <Sidebar />
      <Column
        flex={1}
        overflow="auto"
        height="100%"
        padding={5}
        bg="pageBackground"
        borderRadius={1}
        {...rest}
      >
        {children}
      </Column>
    </Row>
  </Column>
);

export default Layout;
