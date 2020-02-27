import React from 'react';
import styled from 'styled-components';

import { Box } from './box';

const StyledValue = styled(Box).attrs({ as: 'span' })`
  word-wrap: break-word;
  min-height: 19px;
`;

StyledValue.defaultProps = {
  color: 'grays.10',
  fontSize: 2,
};

const Value = ({ children, ...props }) => {
  const value = children || '-';
  return <StyledValue {...props}>{value}</StyledValue>;
};

export default Value;
