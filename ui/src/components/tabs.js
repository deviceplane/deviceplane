import React from 'react';
import styled from 'styled-components';
import { useActive, useLinkProps } from 'react-navi';

import { Row } from './core';
import { Btn } from './core/button';

const TabButton = styled(Btn)`
  border: none;
  padding: 10px 14px;
  cursor: ${props => (props.active ? 'default' : 'pointer')};
  text-transform: uppercase;
  white-space: nowrap;
  text-decoration: none !important;
  box-shadow: none !important;
  border-color: none !important;
  pointer-events: ${props => (props.disabled ? 'none' : 'auto')};
  color: ${props =>
    props.active ? props.theme.colors.primary : props.theme.colors.white};
  background-color: ${props =>
    props.active ? props.theme.colors.black : props.theme.colors.grays[3]};
  &:focus,
  &:hover {
    color: ${props =>
      props.active
        ? props.theme.colors.primary
        : props.theme.colors.pureWhite} !important;
    background-color: ${props =>
      props.active
        ? props.theme.colors.black
        : props.theme.colors.grays[0]} !important;
  }
  margin: 0 12px;
`;

const TabLink = styled(TabButton).attrs({ as: 'a' })`
  text-decoration: none !important;
`;

const Tab = ({ title, tooltip, href, onClick, disabled, active = true }) => {
  if (href) {
    return (
      <TabLink
        {...useLinkProps({ href })}
        active={useActive(href, { exact: false })}
        disabled={disabled}
      >
        {title}
      </TabLink>
    );
  }

  return (
    <TabButton onClick={onClick} active={active} disabled={disabled}>
      {title}
    </TabButton>
  );
};

const Tabs = ({ content = [] }) => {
  return (
    <Row marginX={4}>
      {content.map(tab => (
        <Tab key={tab.title} {...tab} />
      ))}
    </Row>
  );
};

export default Tabs;
