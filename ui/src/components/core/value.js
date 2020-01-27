import React from 'react';
import styled from 'styled-components';
import { space, color, typography } from 'styled-system';

const StyledValue = styled.span`
    word-wrap: break-word;
    min-height: 19px;
    ${color} ${space} ${typography}
`;

StyledValue.defaultProps = {
  color: 'grays.10',
  fontSize: 2,
};

const Value = ({ children }) => {
  const value = children || '-';
  return <StyledValue>{value}</StyledValue>;
};

export default Value;
