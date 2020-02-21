import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints } from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import { Text, Column, Form, toaster } from '../components/core';

const endpoint = endpoints.resetPassword();
const validationSchema = yup.object().shape({
  email: validators.email.required(),
});
const errorMessages = {
  default:
    'There was an error with your e-mail. Please contact us at support@deviceplane.com.',
};

const PasswordReset = () => {
  const navigation = useNavigation();
  return (
    <Column
      flex={1}
      alignItems="center"
      justifyContent="center"
      paddingY={[0, 0, 6]}
      height="100%"
      overflow="auto"
      bg={['black', 'black', 'pageBackground']}
    >
      <Card logo size="medium" title="Reset Password">
        <Text marginBottom={6} fontWeight={1}>
          You will receive an email with a link to reset your password.
        </Text>
        <Form
          endpoint={endpoint}
          onSuccess={() => {
            navigation.navigate(`/login`);
            toaster.success(
              'Password recovery email sent. Please check your email to reset your password.'
            );
          }}
          validationSchema={validationSchema}
          errorMessages={errorMessages}
          submitLabel="Reset Password"
        >
          <Field
            autoFocus
            autoComplete="on"
            type="email"
            label="Email address"
            name="email"
          />
        </Form>
      </Card>
    </Column>
  );
};

export default PasswordReset;
