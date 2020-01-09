import styled from 'styled-components';
import { space, color, typography, border } from 'styled-system';

const Badge = styled.div`
  color: ${props => props.theme.colors.black};
  text-transform: uppercase;
  border-radius: 6px;
  ${color} ${space} ${typography} ${border};
  padding: 2px 4px;
`;

Badge.defaultProps = {
  fontSize: 0,
  fontWeight: 2,
};

export default Badge;
