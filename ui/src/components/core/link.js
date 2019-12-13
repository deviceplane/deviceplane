import React from 'react';
import styled from 'styled-components';
import { space, layout, color, typography } from 'styled-system';
import { useLinkProps } from 'react-navi';

const A = styled.a`
  text-decoration: none;
  cursor: pointer;

  &:hover {
    text-decoration: underline;
  }

  ${color} ${layout} ${space} ${typography}
`;

A.defaultProps = {
  color: 'primary',
};

const Link = ({ children, href, ...rest }) => {
  const linkProps = useLinkProps({ href });
  return (
    <A {...linkProps} {...rest}>
      {children}
    </A>
  );
};

export default Link;
