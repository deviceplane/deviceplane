import React from 'react';
import styled from 'styled-components';
import {
  typography,
  color,
  space,
  layout,
  border,
  shadow,
} from 'styled-system';

import { Row } from './box';
import Text from './text';

const StyledInput = styled.input.attrs(props => ({
  type: 'radio',
}))``;

const Radio = ({ label, checked, onChange, ...props }) => {
  const handleChange = event => onChange(event.target.checked);

  return (
    <Row
      {...props}
      onClick={() => onChange(!checked)}
      style={{ cursor: 'pointer' }}
      alignSelf="flex-start"
    >
      <StyledInput
        checked={checked}
        onChange={handleChange}
        style={{ cursor: 'pointer' }}
      />
      <Text marginLeft={1} style={{ whiteSpace: 'nowrap' }}>
        {label}
      </Text>
    </Row>
  );
};

export default Radio;
