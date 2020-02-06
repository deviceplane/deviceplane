import React from 'react';
import SquareLoader from 'react-spinners/SquareLoader';
import styled, { keyframes } from 'styled-components';

import { Column } from './core';

const show = keyframes`
0% {
  display: none;
  opacity: 0;
}
50% {
  display: flex;
}
100% {
  opacity: 1;
}
`;

const Container = styled(Column)`
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 16px;
  display: none;
  animation: ${show} 500ms 500ms forward-fill;

  & > div {
    border-radius: 3px;
  }
`;

const Spinner = ({ size = 32 }) => (
  <Container>
    <SquareLoader color="white" size={size} />
  </Container>
);

export default Spinner;
