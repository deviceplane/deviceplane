import React from 'react';
import styled from 'styled-components';
import { useActive, useLinkProps } from 'react-navi';

import { Row } from './core';
import { Btn } from './core/button';

const TabButton = styled(Btn)`
  margin: 0;
  padding: 0;
  cursor: ${props => (props.active ? 'default' : 'pointer')};
  text-transform: uppercase;
  white-space: nowrap;
  text-decoration: none !important;
  border: none;
  border-radius: 0;
  pointer-events: ${props => (props.disabled ? 'none' : 'auto')};
  background-color: ${props => props.theme.colors.black};
  color: ${props =>
    props.active ? props.theme.colors.primary : props.theme.colors.grays[13]};

  &:not(:last-child) {
    padding-right: 12px;
    margin-right: 12px;
    border-right: 2px solid ${props => props.theme.colors.grays[5]};
  }

  &:focus,
  &:hover {
    color: ${props =>
      props.active
        ? props.theme.colors.primary
        : props.theme.colors.pureWhite} !important;
  }
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
