import React, { useState, forwardRef } from 'react';
import styled from 'styled-components';
import {
  typography,
  color,
  space,
  layout,
  border,
  shadow,
} from 'styled-system';
import { Icon } from 'evergreen-ui';

import theme from '../../theme';

const Container = styled.div`
  display: flex;
  flex: 1;
  align-items: center;
  position: relative;

  & svg {
    position: absolute;
    right: 12px;
    cursor: pointer;
  }
`;

const StyledInput = styled.input`
  border: 1px solid ${props => props.theme.colors.grays[0]};
  outline: none;
  margin: 0;
  transition: border-color 200ms;
  width: 100%;

  &:focus {
    border-color: ${props => props.theme.colors.primary};
  }

  &::placeholder {
    font-size: 16px;
    color: inherit;
    opacity: .75;
  }

  -webkit-autofill,
  -webkit-autofill:hover, 
  -webkit-autofill:focus, 
  -webkit-autofill:active  {
      box-shadow: 0 0 0 30px ${props =>
        props.theme.colors.grays[1]} inset !important;
  }

  ${space} ${border} ${layout} ${color} ${typography} ${shadow}
`;

StyledInput.defaultProps = {
  color: 'grays.11',
  bg: 'grays.0',
  borderRadius: 1,
  fontWeight: 2,
  boxShadow: 0,
  fontSize: 2,
};

const Input = forwardRef((props, ref) => {
  const [type, setType] = useState(props.type);
  return (
    <Container>
      <StyledInput
        ref={ref}
        padding={3}
        paddingRight={props.type === 'password' ? 6 : 3}
        {...props}
        type={type}
      />
      {props.type === 'password' && (
        <Icon
          size={14}
          icon={type === 'password' ? 'eye-off' : 'eye-open'}
          color={theme.colors.white}
          onClick={() => setType(type === 'password' ? 'text' : 'password')}
        />
      )}
    </Container>
  );
});

export default Input;
