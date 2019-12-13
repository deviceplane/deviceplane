import styled from 'styled-components';
import { space, color, typography } from 'styled-system';

const Badge = styled.div`
  background-color: ${props => props.theme.colors.white};
  color: ${props => props.theme.colors.black};
  text-transform: uppercase;
  border-radius: 2px;
  ${color} ${space} ${typography};
  padding: 4px 6px;
`;

Badge.defaultProps = {
  fontSize: 0,
  fontWeight: 3,
};

export default Badge;
