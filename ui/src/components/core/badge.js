import styled from 'styled-components';
import { space, color, typography, border } from 'styled-system';

const Badge = styled.div`
  color: ${props => props.theme.colors.black};
  text-transform: uppercase;
  border-radius: 3px;
  ${color} ${space} ${typography} ${border};
  padding: 2px 4px;
`;

Badge.defaultProps = {
  fontSize: 0,
  fontWeight: 3,
};

export default Badge;
