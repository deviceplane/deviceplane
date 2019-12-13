import styled from 'styled-components';
import { layout, space, color, typography, border } from 'styled-system';

const Text = styled.span`
  word-wrap: break-word;
  text-overflow: ellipsis;
  ${color} ${space} ${typography} ${layout} ${border}
`;

Text.defaultProps = {
  color: 'white',
};

export default Text;
