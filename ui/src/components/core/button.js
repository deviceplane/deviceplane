import React from 'react';
import styled from 'styled-components';
import { variant } from 'styled-system';
import { useLinkProps } from 'react-navi';

import theme from '../../theme';
import { Box } from './box';

const variants = {
  variants: {
    primary: {
      color: 'black',
      bg: 'primary',
      border: 0,
      opacity: 0.9,
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        bg: 'transparent',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    primaryGray: {
      color: 'black',
      bg: 'grays.3',
      border: 0,
      borderColor: 'grays.3',
      opacity: 0.9,
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        bg: 'transparent',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    secondary: {
      color: 'primary',
      bg: 'transparent',
      border: 0,
      borderColor: 'primary',
      opacity: 0.9,
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    tertiary: {
      color: 'primary',
      bg: 'transparent',
      padding: '5px 6px',
      opacity: 0.9,
      border: 0,

      height: 'min-content',
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    tertiaryGreen: {
      color: 'green',
      bg: 'transparent',
      padding: '5px 6px',
      opacity: 0.9,
      border: 0,
      borderColor: 'green',
      height: 'min-content',
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    tertiaryDanger: {
      color: 'red',
      bg: 'transparent',
      padding: '5px 6px',
      opacity: 0.9,
      height: 'min-content',
      border: 0,
      borderColor: 'red',
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    danger: {
      color: 'red',
      bg: 'transparent',
      border: 0,
      borderColor: 'red',
      opacity: 0.9,
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
        borderColor: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        opacity: 1,
      },
    },
    text: {
      color: 'grays.10',
      bg: 'transparent',
      padding: 0,
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        color: 'pureWhite',
      },
    },
    textPrimary: {
      color: 'primary',
      bg: 'transparent',
      padding: 0,
      '&:disabled, &[disabled]': {
        color: 'pureWhite',
      },
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        color: 'pureWhite',
      },
    },
    icon: {
      bg: 'grays.3',
      width: '32px',
      height: '32px',
      padding: 0,
      border: 0,
      borderColor: 'grays.3',
      borderRadius: '16px',
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        borderColor: 'primary',
      },
      '&:disabled svg': {
        fill: theme.colors.pureWhite,
        stroke: theme.colors.pureWhite,
      },
    },
    iconSecondary: {
      bg: 'grays.3',
      width: '32px',
      height: '32px',
      padding: 0,
      border: 0,
      borderColor: 'grays.3',
      borderRadius: '16px',
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        borderColor: 'pureWhite',
      },
      '&:disabled svg': {
        fill: theme.colors.pureWhite,
        stroke: theme.colors.pureWhite,
      },
    },
    iconDanger: {
      bg: 'grays.3',
      width: '32px',
      height: '32px',
      padding: 0,
      border: 0,
      borderColor: 'grays.3',
      borderRadius: '16px',
      '&:not(:disabled):not([disabled]):hover, &:not(:disabled):not([disabled]):focus': {
        borderColor: 'red',
      },
      '&:disabled svg': {
        fill: theme.colors.pureWhite,
        stroke: theme.colors.pureWhite,
      },
    },
  },
};

export const Btn = styled(Box).attrs({ as: 'button' })`
  display: flex;
  align-items: center;
  justify-content: center;
  appearance: none;
  border: none;
  outline: none;
  font-family: inherit;
  cursor: pointer;
  transition: ${props => props.theme.transitions[0]};
  transform: translateZ(0);
  backface-visibility: hidden;
  white-space: nowrap;
  font-size: 12px;
  padding: 10px 12px;
  text-transform: uppercase;
  text-renderering: geometricPercision;
  border-radius: 2px;
  flex-shrink: 0;

  &:disabled,
  &[disabled] {
    cursor: not-allowed;
    opacity: 0.4;
  }
  &:focus {
    outline: none;
  }

  ${variant(variants)}
`;
Btn.defaultProps = {
  variant: 'primary',
  fontWeight: 2,
};

export const LinkButton = styled(Btn).attrs({
  as: 'a',
})`
  text-decoration: none;
`;

const Button = ({ href, title, onClick, newTab, ...rest }) => {
  if (href) {
    return (
      <LinkButton
        {...(newTab
          ? { href, target: '_blank', rel: 'noopener noreferrer' }
          : useLinkProps({ href, onClick }))}
        {...rest}
      >
        {title}
      </LinkButton>
    );
  }

  return (
    <Btn onClick={onClick} {...rest}>
      {title}
    </Btn>
  );
};

Button.defaultProps = {
  href: null,
  title: '',
};

export default Button;
