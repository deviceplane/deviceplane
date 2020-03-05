import React, { useRef, useEffect } from 'react';
import styled from 'styled-components';
import { space, layout, flexbox } from 'styled-system';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';

import { useMutation } from '../../api';
import Button from './button';
import Alert from '../alert';
import { Row } from './box';

export const StyledForm = styled.form`
  display: flex;
  flex-direction: column;
  margin: 0;
  padding: 0;

  ${flexbox} ${space} ${layout}
`;

export const Form = ({
  onSubmit,
  onData,
  onCancel,
  onSuccess = () => {},
  onError = () => {},
  endpoint = {},
  validationSchema,
  defaultValues,
  errorMessages,
  alert,
  submitLabel = 'Update',
  submitDisabled,
  submitVariant,
  children,
}) => {
  const {
    handleSubmit,
    register,
    control,
    errors,
    formState,
    getValues,
  } = useForm({
    defaultValues,
    validationSchema,
  });
  const submitButtonRef = useRef(null);
  const navigation = useNavigation();
  const [mutate, { data, success, error }] = useMutation(endpoint.url, {
    method: endpoint.method,
    errors: errorMessages,
  });
  useEffect(() => {
    if (success) {
      onSuccess({ ...getValues({ nest: true }), ...data });
    }
  }, [success]);
  useEffect(() => {
    if (error) {
      onError(error);
    }
  }, [error]);

  const submit = data => {
    submitButtonRef.current.blur();
    if (onData) {
      data = onData(data);
    }
    if (onSubmit) {
      return onSubmit(data);
    }
    if (endpoint) {
      mutate(data);
    }
  };

  const modifyChild = child =>
    typeof child === 'object' && child.props.name
      ? React.createElement(child.type, {
          ...{
            ...child.props,
            register,
            control,
            errors: errors[child.props.name],
            key: child.props.name,
          },
        })
      : child;

  return (
    <StyledForm onSubmit={handleSubmit(submit)}>
      <Alert
        show={alert || error}
        variant="error"
        description={alert || error}
      />
      {Array.isArray(children)
        ? children.map(modifyChild)
        : modifyChild(children)}
      <Row>
        <Button
          ref={submitButtonRef}
          variant={submitVariant}
          type="submit"
          title={submitLabel}
          disabled={submitDisabled || formState.isSubmitting}
        />
        {onCancel && (
          <Button
            variant="text"
            title="Cancel"
            href={onCancel}
            marginLeft={5}
          />
        )}
      </Row>
    </StyledForm>
  );
};
