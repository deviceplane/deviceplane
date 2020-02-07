import React from 'react';
import styled from 'styled-components';
import { useActive } from 'react-navi';

import AvatarMenu from './avatar-menu';
import ProjectSelector from './project-selector';
import Logo from './icons/logo';
import { Row, Link, Icon } from './core';

const HeaderLink = styled(Link)`
  display: flex;
  align-items: center;
  color: ${props => props.theme.colors.white};
  font-weight: ${props => props.theme.fontWeights[1]};
  font-size: ${props => props.theme.fontSizes[1]}px;
  &:hover {
    color: ${props => props.theme.colors.primary};
    border-color: ${props => props.theme.colors.primary};
  }

  &:hover svg {
    fill: ${props => props.theme.colors.primary};
  }
`;

const Header = ({ children }) => {
  const isProjectsRoute = useActive('/projects');
  return (
    <Row
      alignItems="center"
      justifyContent="space-between"
      alignSelf="stretch"
      paddingRight={5}
      paddingY={2}
      bg="black"
    >
      <Row
        marginLeft={isProjectsRoute ? 5 : 0}
        width={isProjectsRoute ? 'unset' : '136px'}
        justifyContent="center"
      >
        <Logo />
      </Row>
      <Row flex={1} alignItems="center">
        <ProjectSelector />
      </Row>
      <Row justifyContent="center" flex={1} marginX={5}>
        {children}
      </Row>
      <Row justifyContent="flex-end" alignItems="center" flex={1}>
        <HeaderLink
          newTab
          href="https://deviceplane.com/docs"
          marginRight={5}
          borderRadius={1}
          padding={1}
          paddingX="6px"
          border={0}
          borderColor="white"
        >
          <Icon icon="manual" size={12} color="white" marginRight={2} />
          Docs
        </HeaderLink>
        <AvatarMenu />
      </Row>
    </Row>
  );
};

export default Header;
