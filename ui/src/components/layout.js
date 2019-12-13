import React from 'react';

import { Column, Row } from './core';
import Header from './header';
import Sidebar from './sidebar';

const Layout = ({ children, header, ...rest }) => (
  <Row minHeight="100%">
    <Sidebar />
    <Column flex={1}>
      <Header>{header}</Header>
      <Column flex={1} {...rest} padding={5}>
        {children}
      </Column>
    </Column>
  </Row>
);

export default Layout;
