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
import { Column, Radio } from './core';

const RadioGroup = ({ value, options = [], onChange }) => {
  return (
    <Column>
      {options.map(option => (
        <Radio
          onChange={() => onChange(option.value)}
          key={option.value}
          label={option.label}
          marginBottom={2}
          checked={option.value === value}
        />
      ))}
    </Column>
  );
};

export default RadioGroup;
