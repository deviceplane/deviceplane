import React from 'react';

import { Box, Column, Row } from './core';
import Header from './header';
import Sidebar from './sidebar';

const Layout = ({ children, header, ...rest }) => (
  <Box>
    <Header>{header}</Header>
    <Row>
      <Sidebar />
      <Column
        height="calc(100vh - 64px)"
        flex={1}
        alignSelf="stretch"
        overflow="auto"
        bg="pageBackground"
        borderRadius={1}
        padding={5}
        {...rest}
      >
        {children}
      </Column>
    </Row>
  </Box>
);

export default Layout;
