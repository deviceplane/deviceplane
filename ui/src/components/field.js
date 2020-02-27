import React, { forwardRef, useState } from 'react';
import styled from 'styled-components';
import { space, color, typography } from 'styled-system';
import { Controller } from 'react-hook-form';

import utils from '../utils';
import Editor from './editor';
import {
  Group,
  Row,
  Column,
  Input,
  Textarea,
  Label,
  Text,
  Icon,
  Select,
  MultiSelect,
  Checkbox,
  RadioGroup,
} from './core';

const PasswordButton = styled(Row)`
  user-select: none;
  cursor: pointer;
  align-items: center;
  transition: ${props => props.theme.transitions[0]};
  &:hover span {
    color: ${props => props.theme.colors.pureWhite};
  }
  &:hover svg {
    fill: ${props => props.theme.colors.pureWhite};
  }
`;

const FieldLabel = styled.label`
${space} ${color} ${typography}
`;
FieldLabel.defaultProps = {
  ...Label.defaultProps,
  marginBottom: 0,
};

const Field = forwardRef(
  (
    {
      label,
      hint,
      description,
      name,
      register,
      control,
      autoComplete = 'off',
      multi,
      inline,
      flex,
      errors = [],
      marginBottom = Group.defaultProps.marginBottom,
      ...props
    },
    ref
  ) => {
    const [type, setType] = useState(props.type);

    const getComponent = () => {
      switch (type) {
        case 'editor':
          return (
            <Controller
              name={name}
              id={name}
              control={control}
              as={<Editor />}
              {...props}
            />
          );
        case 'multiselect':
          return (
            <Controller
              multi
              name={name}
              id={name}
              control={control}
              onChange={([selected]) => ({ value: selected })}
              as={<MultiSelect />}
              {...props}
            />
          );
        case 'checkbox':
          return (
            <Controller
              name={name}
              id={name}
              control={control}
              label={label}
              as={<Checkbox />}
              {...props}
            />
          );
        case 'radiogroup':
          return (
            <Controller
              name={name}
              id={name}
              control={control}
              label={label}
              as={<RadioGroup />}
              {...props}
            />
          );
        case 'textarea':
          return <Textarea name={name} id={name} ref={ref} {...props} />;
        case 'select':
          return <Select name={name} id={name} ref={ref} {...props} />;
        default:
          return (
            <Input
              autoComplete={autoComplete}
              name={name}
              id={name}
              ref={ref}
              {...props}
              type={type}
            />
          );
      }
    };

    errors = Array.isArray(errors) ? errors : [errors];

    const component = getComponent();

    if (type === 'checkbox') {
      // allows checkbox to render its own label
      label = null;
    }

    return (
      <Column flex={flex} marginBottom={inline ? 0 : multi ? 4 : marginBottom}>
        {(label || description) && (
          <Column marginBottom={Label.defaultProps.marginBottom}>
            <Row justifyContent="space-between">
              {label && <FieldLabel htmlFor={name}>{label}</FieldLabel>}
              {props.type === 'password' && (
                <PasswordButton
                  onClick={() =>
                    setType(type => (type === 'password' ? 'text' : 'password'))
                  }
                >
                  <Icon
                    icon={type === 'password' ? 'eye-open' : 'eye-off'}
                    size={16}
                    color="primary"
                  />
                  <Text
                    fontSize={0}
                    fontWeight={2}
                    width="40px"
                    textAlign="right"
                    color="primary"
                  >
                    {type === 'password' ? 'SHOW' : 'HIDE'}
                  </Text>
                </PasswordButton>
              )}
            </Row>

            {description && (
              <Text marginTop={2} fontSize={1} fontWeight={1} color="grays.8">
                {description}
              </Text>
            )}
          </Column>
        )}
        {component}
        {hint && (
          <Text marginTop={2} fontSize={0} fontWeight={1} color="grays.8">
            {hint}
          </Text>
        )}
        {errors.map(({ message }) => (
          <Row marginTop={2} alignItems="flex-start" textAlign="left">
            <Icon
              icon="warning-sign"
              color="red"
              size={16}
              flexShrink={0}
              marginTop="1px"
            />
            <Text color="red" marginLeft={2}>
              {utils.capitalize(message)}
            </Text>
          </Row>
        ))}
      </Column>
    );
  }
);

export default Field;
