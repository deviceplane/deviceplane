import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';

import api from '../api';
import validators from '../validators';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Text, Column, Row, Form, Button, toaster } from '../components/core';
import * as auth0 from '../lib/auth0';

const validationSchema = yup.object().shape({
  email: validators.email.required(),
  password: yup.string().required(),
});

const SsoCallback = ({
  route: {
    data: { params, context },
  },
}) => {
  const { register, handleSubmit, errors } = useForm({
    validationSchema,
  });
  const navigation = useNavigation();
  const [status, setStatus] = useState('Loading...');
  const [firstRender, setFirstRender] = useState(true);

  const submit = async data => {
    try {
      await api.loginSSO(data);
      const response = await api.user();
      context.setCurrentUser(response.data);
      navigation.navigate(
        params.redirectTo ? decodeURIComponent(params.redirectTo) : '/projects'
      );
    } catch (error) {
      var msg = 'Error authenticating: ' + JSON.stringify(error);
      setStatus(msg);
      toaster.danger(msg);
    }
  };

  if (firstRender) {
    if (window.location.hash) {
      var vars = window.location.hash
        .substring(1)
        .split('&')
        .map(x => x.split('='))
        .reduce((map, x) => {
          map[x[0]] = decodeURIComponent(x[1]);
          return map;
        }, {});
      navigation.navigate('#');
      submit(vars);
    } else {
      setStatus('Error: SSO data not returned.');
    }
    setFirstRender(false);
  }

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
      <Text color={'white'} fontSize={4} fontWeight={2} marginLeft={2}>
        {status}
      </Text>

      <Row marginTop={5}>
        <Button justifyContent="center" title="Go back" href="/login" />
      </Row>
    </Column>
  );
};

export default SsoCallback;
