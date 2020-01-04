import React, { forwardRef, useState } from 'react';
import { RHFInput } from 'react-hook-form-input';
import styled from 'styled-components';
import { space, color, typography } from 'styled-system';
import { Icon } from 'evergreen-ui';

import theme from '../theme';
import utils from '../utils';
import { Group, Row, Column, Input, Textarea, Label, Text } from './core';

const Container = styled(Group)`
  margin-bottom: ${props =>
    props.multi
      ? props.theme.sizes[2]
      : props.theme.sizes[Group.defaultProps.marginBottom]}px;

  &:last-of-type {
    margin-bottom: ${props =>
      props.theme.sizes[Group.defaultProps.marginBottom]}px !important;
  }
`;

const PasswordButton = styled(Row)`
  user-select: none;
  cursor: pointer;
  align-items: center;

  & span {
    transition: color 150ms;
  }
  &:hover span {
    color: ${props => props.theme.colors.primary};
  }
  & svg {
    transition: fill 150ms;
  }
  &:hover svg {
    fill: ${props => props.theme.colors.primary} !important;
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
      as,
      setValue,
      register,
      onChangeEvent,
      autoComplete = 'off',
      multi,
      errors = [],
      ...props
    },
    ref
  ) => {
    const [type, setType] = useState(props.type);

    errors = Array.isArray(errors) ? errors : [errors];
    const getComponent = () => {
      if (as) {
        return (
          <RHFInput
            as={as}
            id={name}
            name={name}
            register={register}
            setValue={setValue}
            onChangeEvent={data => ({ value: data[0] })}
            {...props}
          />
        );
      }

      switch (type) {
        case 'textarea':
          return <Textarea name={name} id={name} ref={ref} {...props} />;
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

    return (
      <Container multi={multi}>
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
                  />
                  <Text
                    fontSize={0}
                    marginLeft={2}
                    fontWeight={3}
                    width="40px"
                    textAlign="right"
                  >
                    {type === 'password' ? 'SHOW' : 'HIDE'}
                  </Text>
                </PasswordButton>
              )}
            </Row>

            {description && (
              <Text marginTop={2} fontSize={1} color="grays.8">
                {description}
              </Text>
            )}
          </Column>
        )}
        <Row>{getComponent()}</Row>
        {hint && (
          <Text marginTop={2} fontSize={0} color="grays.8">
            {hint}
          </Text>
        )}
        {errors.map(({ message }) => (
          <Row marginTop={2} alignItems="flex-start">
            <Icon
              icon="error"
              color={theme.colors.red}
              size={14}
              flexShrink={0}
              marginTop={2}
            />
            <Text color="red" marginLeft={2}>
              {utils.capitalize(message)}
            </Text>
          </Row>
        ))}
      </Container>
    );
  }
);

export default Field;
