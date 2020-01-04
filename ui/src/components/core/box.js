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
} from 'styled-system';

export const Box = styled.div`
  ${space} ${layout} ${color} ${border} ${typography} ${shadow} ${position}
`;

export const Row = styled(Box)`
  ${flexbox}

  display: flex;
`;

export const Column = styled(Row)`
  flex-direction: column;
`;
