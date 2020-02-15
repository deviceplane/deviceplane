import styled from 'styled-components';
import {
  space,
  layout,
  color,
  border,
  flexbox,
  typography,
  shadow,
  position,
  grid,
} from 'styled-system';

export const Box = styled.div`
  ${space} ${layout} ${color} ${border} ${typography} ${shadow} ${position} ${flexbox} ${grid}
`;

export const Row = styled(Box)`
  display: flex;
`;

export const Column = styled(Row)`
  flex-direction: column;
`;

export const Grid = styled(Box)`
  display: grid;
`;
