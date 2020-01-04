import React from 'react';
import styled from 'styled-components';
import { space, color, typography } from 'styled-system';
import { useActive, useLinkProps } from 'react-navi';

import { Row } from './core';

const Container = styled(Row)``;

const styles = `
    border: none;
    outline: none;
    border-radius: 0;
    transition: background-color 150ms;
    border-radius: 4px;
    padding: 10px 14px;
    user-select: none;
    cursor: pointer;
    text-transform: uppercase;
    white-space: nowrap;

    &:not(:last-child) {
        margin-right: 18px;
    }
`;

const LinkTab = styled.a`
  text-decoration: none !important;

  ${color} ${typography} ${space}

  ${styles}

  pointer-events: ${props => (props.disabled ? 'none' : 'auto')};
  font-size: ${props => props.theme.fontSizes[1]}px;
  font-weight: ${props => props.theme.fontWeights[2]};
  color: ${props =>
    props.active ? props.theme.colors.primary : props.theme.colors.white};
  background-color: ${props =>
    props.active ? props.theme.colors.black : 'transparent'};

  opacity: ${props => (props.disabled ? 0.4 : 1)};

  &:hover {
    background-color: ${props =>
      props.active
        ? props.theme.colors.black
        : props.disabled
        ? 'transparent'
        : props.theme.colors.grays[1]};
  }
`;

const ButtonTab = styled.button`
  appearance: none;

  ${color} ${typography} ${space}

  ${styles}

  pointer-events: ${props => (props.disabled ? 'none' : 'auto')};
  font-size: ${props => props.theme.fontSizes[1]}px;
  font-weight: ${props => props.theme.fontWeights[2]};
  color: ${props =>
    props.active ? props.theme.colors.primary : props.theme.colors.white};
  background-color: ${props =>
    props.active ? props.theme.colors.black : 'transparent'};

  opacity: ${props => (props.disabled ? 0.4 : 1)};

  &:hover {
    background-color: ${props =>
      props.active
        ? props.theme.colors.black
        : props.disabled
        ? 'transparent'
        : props.theme.colors.grays[1]};
  }
`;

const Tab = ({ title, tooltip, href, onClick, disabled, active = true }) => {
  if (href) {
    return (
      <LinkTab
        {...useLinkProps({ href })}
        active={useActive(href, { exact: false })}
        disabled={disabled}
      >
        {title}
      </LinkTab>
    );
  }

  return (
    <ButtonTab onClick={onClick} active={active} disabled={disabled}>
      {title}
    </ButtonTab>
  );
};

const Tabs = ({ content = [] }) => {
  return (
    <Container marginX={4}>
      {content.map(tab => (
        <Tab key={tab.title} {...tab} />
      ))}
    </Container>
  );
};

export default Tabs;
