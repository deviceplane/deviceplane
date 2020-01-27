import styled from 'styled-components';
import {
  typography,
  color,
  space,
  layout,
  border,
  shadow,
  variant,
} from 'styled-system';

const Input = styled.input.attrs(props => ({
  spellCheck: props.spellCheck || 'off',
  autoCorrect: props.autoCorrect || 'off',
  autoComplete: props.autoComplete || 'off',
}))`
  border: 1px solid ${props => props.theme.colors.white};
  outline: none;
  margin: 0;
  transition: ${props => props.theme.transitions[0]};
  width: 100%;
  padding: 8px;
  caret-color: ${props => props.theme.colors.primary};
  font-weight: 300;

  &:focus {
    border-color: ${props => props.theme.colors.primary};
  }

  font-size: 16px;

  &::placeholder {
    font-size: 16px;
    font-weight: 400;
    color: ${props => props.theme.colors.grays[9]};
  }

  -webkit-autofill,
  -webkit-autofill:hover, 
  -webkit-autofill:focus, 
  -webkit-autofill:active  {
      box-shadow: 0 0 0 30px ${props =>
        props.theme.colors.grays[1]} inset !important;
  }

  ${variant({
    variants: {
      small: {
        padding: '4px',
        fontSize: 1,
        '::placeholder': {
          fontSize: 1,
        },
      },
    },
  })}

  ${space} ${border} ${layout} ${color} ${typography} ${shadow}
`;

Input.defaultProps = {
  color: 'white',
  bg: 'grays.0',
  borderRadius: 1,
};

export default Input;
