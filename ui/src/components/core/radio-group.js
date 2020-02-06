import React from 'react';

import { Column } from './box';
import Radio from './radio';

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
