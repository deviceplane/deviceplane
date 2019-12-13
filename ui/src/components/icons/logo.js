import React from 'react';
import styled from 'styled-components';

import { Row } from '../core';

const Container = styled(Row)`
  @keyframes enter {
    to {
      transform: rotate(360deg) scale(1.25);
    }
  }

  &:hover {
    animation: enter 2s ease-in-out 5s alternate infinite;
  }
`;

const Logo = ({ size = 45, color = 'white' }) => (
  <Container
    width={size}
    height={size}
    alignItems="center"
    justifyContent="center"
    overflow="hidden"
  >
    <svg
      width={175}
      height={152}
      viewBox="0 0 175 152"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M119.2 39.7L159.2 108.3L134.1 151.8L94.3 83.5999H40L65.4 39.7H119.2Z"
        fill={color}
      />
      <path
        d="M127.8 145.6H46.6L6 75.3L46.6 5H127.8L167.4 73.6"
        stroke={color}
        strokeWidth="10"
        strokeMiterlimit="3"
        strokeLinecap="square"
      />
    </svg>
  </Container>
);

export default Logo;
