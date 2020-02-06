import styled from 'styled-components';

import { Box } from './box';

const MenuItem = styled(Box).attrs({ as: 'button' })`
  background: none;
  appearance: none;
  cursor: pointer;
  border: none;
  text-align: left;
  text-transform: uppercase;

  &:hover {
    color: ${props => props.theme.colors.pureWhite};
    background-color: ${props => props.theme.colors.grays[3]};
  }
`;

MenuItem.defaultProps = {
  paddingY: 1,
  color: 'white',
  fontSize: 1,
  fontWeight: 2,
  paddingX: 3,
};

export default MenuItem;
