import React from 'react';
import styled from 'styled-components';
import { space } from 'styled-system';

const HiddenCheckbox = styled.input.attrs({ type: 'checkbox' })`
  // Hide checkbox visually but remain accessible to screen readers.
  // Source: https://polished.js.org/docs/#hidevisually
  border: 0;
  clip: rect(0 0 0 0);
  clippath: inset(50%);
  height: 1px;
  margin: -1px;
  overflow: hidden;
  padding: 0;
  position: absolute;
  white-space: nowrap;
  width: 1px;
`;

const Icon = styled.svg`
  fill: none;
  stroke: black;
  stroke-width: 3px;
`;

const StyledCheckbox = styled.div`
  display: inline-block;
  width: 18px;
  height: 18px;
  background: ${props =>
    props.checked ? props.theme.colors.primary : props.theme.colors.white};
  border-radius: 2px;
  transition: ${props => props.theme.transitions[0]};
  cursor: pointer;

  ${Icon} {
    visibility: ${props => (props.checked ? 'visible' : 'hidden')};
  }
`;

const Container = styled.div`
  display: flex;
  align-items: center;
  user-select: none;

  ${space}
`;

const Text = styled.label`
  cursor: pointer;
  ${space}
`;

const Checkbox = ({ value, label, onChange, ...props }) => {
  return (
    <Container onClick={() => onChange(!value)}>
      <HiddenCheckbox checked={value} id={label} {...props} readOnly />
      <StyledCheckbox checked={value}>
        <Icon viewBox="0 0 24 24">
          <polyline points="20 6 9 17 4 12" />
        </Icon>
      </StyledCheckbox>
      <Text paddingLeft={2} htmlFor={label}>
        {label}
      </Text>
    </Container>
  );
};

export default Checkbox;
