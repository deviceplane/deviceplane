import styled from 'styled-components';
import { space, color, typography } from 'styled-system';

const Label = styled.span`
${space} ${color} ${typography}
`;

Label.defaultProps = {
  color: 'white',
  fontWeight: 1,
  fontSize: 2,
  marginBottom: 2,
};

export default Label;
