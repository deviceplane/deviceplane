import React from 'react';
import { IconSvgPaths20 } from '@blueprintjs/icons';
import styled from 'styled-components';
import { space, layout, border, position } from 'styled-system';

import theme from '../../theme';

const Svg = styled.svg`
flex-shrink: 0;
  ${space} ${layout} ${border} ${position}
`;

const Icon = ({ color, icon, alt, size = 16, ...props }) => {
  const paths = IconSvgPaths20;

  color = theme.colors[color];

  return (
    <Svg
      alt={alt}
      fill={color}
      width={`${size}px`}
      height={`${size}px`}
      viewBox="0 0 20 20"
      {...props}
    >
      {paths[icon].map((d, i) => (
        <path key={i} d={d} fillRule="evenodd" />
      ))}
    </Svg>
  );
};

export default Icon;
