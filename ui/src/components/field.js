import React, { forwardRef } from 'react';
import { RHFInput } from 'react-hook-form-input';
import styled from 'styled-components';
import { space, color, typography } from 'styled-system';
import { Icon } from 'evergreen-ui';

import theme from '../theme';
import utils from '../utils';
import { Group, Row, Column, Input, Textarea, Label, Text } from './core';

const Container = styled(Group)`
  margin-bottom: ${props =>
    props.group
      ? props.theme.sizes[2]
      : props.theme.sizes[Group.defaultProps.marginBottom]}px;

  &:last-of-type {
    margin-bottom: ${props =>
      props.theme.sizes[Group.defaultProps.marginBottom]}px !important;
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
      type,
      name,
      as,
      setValue,
      register,
      onChangeEvent,
      autoComplete = 'off',
      group,
      errors = [],
      ...props
    },
    ref
  ) => {
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
              type={type}
              name={name}
              id={name}
              ref={ref}
              {...props}
            />
          );
      }
    };

    return (
      <Container group={group}>
        {(label || description) && (
          <Column marginBottom={Label.defaultProps.marginBottom}>
            {label && <FieldLabel htmlFor={name}>{label}</FieldLabel>}
            {description && (
              <Text marginTop={2} fontSize={1} color="grays.8">
                {description}
              </Text>
            )}
          </Column>
        )}
        {getComponent()}
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
