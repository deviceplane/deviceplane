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
      var resp;
      switch (params.redirectType) {
        case 'login':
          resp = await api.loginSSO(data);
          break;
        case 'signup':
          resp = await api.signupSSO(data);
          break;
        default:
          var msg = 'Error, invalid redirect type';
          setStatus(msg);
          toaster.danger(msg);
          return;
      }

      const response = await api.user();
      context.setCurrentUser(response.data);
      navigation.navigate(
        params.redirectTo ? decodeURIComponent(params.redirectTo) : '/projects'
      );
    } catch (error) {
      var statusCode = error.response.status;
      var statusText = error.response.statusText;
      if (statusCode == 404) {
        statusText = 'User not found. Have you already signed up?';
      }

      // TODO: just set the toast and navigate back to the page, instead of
      // showing the message and allowing the user to manually use the back button

      var msg = `Error code ${statusCode}: ${statusText}`;
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
      <Text
        color={'white'}
        fontSize={6}
        fontWeight={2}
        marginLeft={2}
        marginBottom={2}
      >
        {params.redirectType.toUpperCase()}
      </Text>
      <Text color={'white'} fontSize={4} fontWeight={2} marginLeft={2}>
        {status}
      </Text>

      <Row marginTop={8}>
        <Button
          justifyContent="center"
          title="Go back"
          onClick={() => {
            navigation.navigate('/' + params.redirectType);
          }}
        />
      </Row>
    </Column>
  );
};

export default SsoCallback;
