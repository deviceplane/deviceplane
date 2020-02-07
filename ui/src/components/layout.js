import React from 'react';

import { Row, Column } from './core';
import Header from './header';
import Sidebar from './sidebar';

const Layout = ({ children, header, ...rest }) => (
  <Row flex={1} flexDirection={['column-reverse', 'column-reverse', 'row']}>
    <Sidebar />
    <Column flex={1}>
      <Header>{header}</Header>

      <Column
        flex={1}
        overflow="auto"
        bg={['black', 'black', 'pageBackground']}
        borderRadius={1}
        padding={5}
        {...rest}
      >
        {children}
      </Column>
    </Column>
  </Row>
);

export default Layout;
