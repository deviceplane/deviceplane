import styled from 'styled-components';
import { space, layout, flexbox } from 'styled-system';

export const Form = styled.form`
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;

  ${flexbox} ${space} ${layout}
`;
