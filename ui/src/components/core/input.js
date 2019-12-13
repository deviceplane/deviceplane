import styled from 'styled-components';
import {
  typography,
  color,
  space,
  layout,
  border,
  shadow,
} from 'styled-system';

const Input = styled.input`
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

Input.defaultProps = {
  color: 'grays.11',
  bg: 'grays.0',
  borderRadius: 1,
  fontWeight: 2,
  boxShadow: 0,
  fontSize: 2,
  padding: 3,
};

export default Input;
