import React from 'react';
import styled from 'styled-components';
import { typography, color, space, layout } from 'styled-system';
import titlify from 'title';

const defaultProps = {
  color: 'white',
  fontWeight: 3,
  margin: 0,
};

const H1 = styled.h1`
  word-break: break-word;
  text-transform: none;
  ${typography} ${color} ${space} ${layout} 
`;

H1.defaultProps = {
  ...defaultProps,
  fontSize: 7,
};

const H2 = styled.h2`
  margin: 0;
  word-break: break-word;
  ${typography} ${color} ${space} ${layout} 
`;

H2.defaultProps = {
  ...defaultProps,
  fontSize: 6,
};

const H3 = styled.h3`
  word-break: break-word;
  ${typography} ${color} ${space} ${layout} 
`;

H3.defaultProps = {
  ...defaultProps,
  fontSize: 5,
};

const Heading = ({ variant, ...rest }) => {
  rest.children = titlify(rest.children);

  switch (variant) {
    case 'secondary':
      return <H2 {...rest} />;
    case 'tertiary':
      return <H3 {...rest} />;
    case 'primary':
    default:
      return <H1 {...rest} />;
  }
};

export default Heading;
