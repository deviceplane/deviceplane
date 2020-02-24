import styled from 'styled-components';

import { Box } from './box';

const Group = styled(Box)`
  display: flex;
  flex-direction: column;

  &:last-child {
    margin-bottom: 0;
  }
`;

Group.defaultProps = {
  marginBottom: 5,
};

export default Group;
