import React from 'react';
import styled from 'styled-components';

import Icon from './icon';
import { Box } from './box';

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

const CheckIcon = styled.svg`
  fill: none;
  stroke: ${props => props.theme.colors.primary};
  stroke-width: 3px;
`;

const StyledCheckbox = styled.div`
  display: inline-block;
  width: 18px;
  height: 18px;
  background: ${props => props.theme.colors.grays[0]};
  border: 1px solid
    ${props =>
      props.checked ? props.theme.colors.primary : props.theme.colors.white};
  border-radius: 2px;
  transition: ${props => props.theme.transitions[0]};
  cursor: pointer;

  &:hover {
    border-color: ${props => props.theme.colors.primary};
  }

  ${CheckIcon} {
    fill: none !important;
    visibility: ${props => (props.checked ? 'visible' : 'hidden')};
  }
`;

const Container = styled(Box)`
  display: flex;
  align-items: center;
  user-select: none;
  align-self: flex-start;

  pointer-events: ${props => (props.readOnly ? 'none' : 'unset')};
`;

const Text = styled(Box).attrs({ as: 'label' })`
  cursor: pointer;
`;

const Checkbox = ({ readOnly, checked, label, onChange, ...props }) => {
  if (readOnly) {
    if (checked) {
      return (
        <Container>
          <CheckIcon viewBox="0 0 24 24" width={20} height={20}>
            <polyline points="20 6 9 17 4 12" />
          </CheckIcon>
        </Container>
      );
    }
    return (
      <Container>
        <Icon icon="cross" size={20} color="grays.10" />
      </Container>
    );
  }
  return (
    <Container onClick={() => onChange(!checked)}>
      <HiddenCheckbox id={label} {...props} readOnly />
      <StyledCheckbox checked={checked}>
        <CheckIcon viewBox="0 0 24 24">
          <polyline points="20 6 9 17 4 12" />
        </CheckIcon>
      </StyledCheckbox>
      {label && (
        <Text paddingLeft={2} htmlFor={label}>
          {label}
        </Text>
      )}
    </Container>
  );
};

export default Checkbox;
