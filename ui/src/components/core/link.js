import React from 'react';
import styled from 'styled-components';
import { space, layout, color, typography } from 'styled-system';
import { useLinkProps } from 'react-navi';

import { Box } from './box';

const A = styled(Box).attrs({ as: 'a' })`
  text-decoration: none;
  cursor: pointer;
  color: ${props => props.theme.colors.primary};

  &:hover {
    color: ${props => props.theme.colors.white};
  }

  ${color} ${layout} ${space} ${typography}
`;

A.defaultProps = {
  fontWeight: 2,
};

const Link = ({ children, href, newTab, ...rest }) => {
  const linkProps = useLinkProps({ href });
  return (
    <A
      {...linkProps}
      {...(newTab ? { target: '_blank', rel: 'noopener noreferrer' } : {})}
      {...rest}
    >
      {children}
    </A>
  );
};

export default Link;
