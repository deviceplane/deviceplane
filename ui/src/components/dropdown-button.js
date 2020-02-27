import React from 'react';
import styled from 'styled-components';

import Popover from './popover';
import { Icon, Text, Box, Row, Link } from './core';

export const Button = styled.button`
  display: flex;
  align-items: center;
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
  font-weight: 500;
  padding: 10px 12px;
  text-transform: uppercase;
  text-renderering: geometricPercision;
  flex: 0;
  background-color: ${props => props.theme.colors.white};
  color: ${props => props.theme.colors.black};
  border-radius: 2px;
  pointer-events: ${props => (props.disabled ? 'none' : 'auto')};

  &:disabled {
    cursor: not-allowed;
    opacity: 0.4;
  }

  &:focus {
    outline: none;
  }
`;

const MenuItem = styled(Box).attrs({ as: 'button' })`
  display: flex;
  cursor: ${props => (props.disabled ? 'not-allowed' : 'pointer')};
  white-space: nowrap;
  border: none;
  outline: none;
  appearance: none;
  background: ${props => props.theme.colors.black};

  opacity: ${props => (props.disabled ? 0.4 : 1)};

  &:hover {
    background-color: ${props => props.theme.colors.grays[3]};
  }
`;
MenuItem.defaultProps = {
  paddingX: '6px',
  paddingY: 2,
  color: 'white',
  fontSize: 1,
  fontWeight: 2,
};

const DropdownButton = ({
  title,
  button,
  content = [],
  disabled,
  width = 'auto',
  top = '43px',
  right,
}) => {
  return (
    <Popover
      disabled={disabled}
      button={({ show }) =>
        button || (
          <Button disabled={disabled}>
            <Text fontSize={0} fontWeight={3} color="inherit">
              {title}
            </Text>
            <Icon icon="caret-down" size={16} color="black" marginLeft={2} />
          </Button>
        )
      }
      content={({ close }) =>
        content.map(({ label, onClick, color, disabled }) => (
          <MenuItem
            disabled={disabled}
            color={color}
            onClick={() => {
              onClick();
              close();
            }}
          >
            {label}
          </MenuItem>
        ))
      }
      top={top}
      right={right}
      width={width}
    ></Popover>
  );
};

export default DropdownButton;
