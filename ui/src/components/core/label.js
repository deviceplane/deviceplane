import styled from 'styled-components';

import { Box } from './box';

const Label = styled(Box).attrs({ as: 'span' })``;

Label.defaultProps = {
  color: 'white',
  fontWeight: 1,
  fontSize: 2,
  marginBottom: 2,
};

export default Label;
