import React from 'react';
import styled from 'styled-components';
import { useActive, useCurrentRoute } from 'react-navi';

import Logo from './icons/logo';
import { Row, Column, Link, Text, Icon } from './core';

import storage from '../storage';

let links;
if (storage.get('legacy') || false) {
  links = [
    {
      title: 'Devices',
      icon: 'multi-select',
      to: '/devices',
    },
    {
      title: 'Provisioning',
      icon: 'projects',
      to: '/provisioning',
    },
    {
      title: 'Applications',
      icon: 'applications',
      to: '/applications',
    },
    {
      title: 'IAM',
      icon: 'people',
      to: '/iam',
    },
    {
      title: 'Settings',
      icon: 'settings',
      to: '/settings',
    },
  ];
} else {
  links = [
    {
      title: 'Devices',
      icon: 'multi-select',
      to: '/devices',
    },
    {
      title: 'Provisioning',
      icon: 'projects',
      to: '/provisioning',
    },
    {
      title: 'IAM',
      icon: 'people',
      to: '/iam',
    },
    {
      title: 'Settings',
      icon: 'settings',
      to: '/settings',
    },
  ];
}

const SidebarLink = styled(Link)`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  transition: ${props => props.theme.transitions[0]};
  text-transform: uppercase;
  text-decoration: none !important;
  color: ${props =>
    props.active ? props.theme.colors.primary : props.theme.colors.grays[12]};
  &:hover {
    color: ${props =>
      props.active ? props.theme.colors.primary : props.theme.colors.pureWhite};
  }
  &:hover svg {
    fill: ${props =>
      props.active ? props.theme.colors.primary : props.theme.colors.pureWhite};
  }
  & svg {
    fill: ${props =>
      props.active
        ? props.theme.colors.primary
        : props.theme.colors.grays[12]}};
  }
  &:last-child {
    margin-top: auto;
  }
`;

const Sidebar = () => {
  const route = useCurrentRoute();

  if (!route) {
    return null;
  }

  const projectSelected = !!route.data.params.project;

  return (
    projectSelected && (
      <Column
        flexDirection={['row', 'row', 'column']}
        bg={['grays.0', 'grays.0', 'black']}
        alignItems={['unset', 'unset', 'stretch']}
        justifyContent={['space-between', 'space-between', 'unset']}
        flexShrink={0}
        overflow="auto"
      >
        <Row
          paddingY={4}
          justifyContent="center"
          display={['none', 'none', 'flex']}
        >
          <Logo size={40} />
        </Row>

        {links.map(({ to, title, icon }) => {
          const href = `/${route.data.params.project}${to}`;

          return (
            <SidebarLink
              href={href}
              paddingY={[2, 2, 4, 4]}
              paddingX={[5, 5, 5, 4]}
              key={title}
              active={useActive(href, { exact: false })}
            >
              <Icon icon={icon} color="white" size={24} />
              <Text
                marginTop={[2, 2, 2, 3]}
                display={['none', 'none', 'none', 'block']}
                fontSize={0}
                color="inherit"
              >
                {title}
              </Text>
            </SidebarLink>
          );
        })}
      </Column>
    )
  );
};

export default Sidebar;
