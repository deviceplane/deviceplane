import styled from 'styled-components';
import { variant } from 'styled-system';
import React, { forwardRef } from 'react';

import { Row, Box } from './box';
import Icon from './icon';

const StyledSelect = styled(Box).attrs({ as: 'select' })`
  background: ${props => props.theme.colors.grays[0]};
  color: ${props => props.theme.colors.white};
  border-radius: ${props => props.theme.radii[1]}px;
  appearance: none;
  padding: 8px;
  font-size: 16px;
  font-weight: ${props => (props.value ? 300 : 400)};
  display: flex;
  flex: 1;
  border: 1px solid ${props => props.theme.colors.white};
  outline: none;

  transition: ${props => props.theme.transition};

  &:focus {
    border-color: ${props => props.theme.colors.primary};
  }

  &:invalid {
    color: ${props => props.theme.colors.grays[9]};
  }

  ${variant({
    variants: {
      small: {
        padding: '4px 6px',
        fontSize: 1,
      },
    },
  })}
`;

const Select = forwardRef(
  (
    {
      name,
      id,
      options,
      disabled,
      autoFocus,
      value,
      required,
      placeholder,
      none = 'There are no options',
      onChange,
      variant,
      ...props
    },
    ref
  ) => {
    return (
      <Row
        flex={1}
        position="relative"
        style={{ cursor: 'pointer' }}
        {...props}
      >
        <StyledSelect
          name={name}
          id={id}
          variant={variant}
          disabled={disabled}
          autoFocus={autoFocus}
          required={required}
          value={value}
          onChange={onChange}
          ref={ref}
        >
          {options.length === 0 && (
            <option value="" disabled selected hidden>
              {none}
            </option>
          )}
          {options.length > 0 && placeholder && (
            <option value="" disabled selected hidden>
              {placeholder}
            </option>
          )}
          {options.map(({ label, value }) => (
            <option value={value}>{label}</option>
          ))}
        </StyledSelect>
        <Icon
          icon="caret-down"
          color="white"
          size={16}
          position="absolute"
          right={2}
          top="25%"
          pointerEvents="none"
        />
      </Row>
    );
  }
);

export default Select;
