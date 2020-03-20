import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import api from '../api';
import validators from '../validators';
import { Text, Column, Row, Button, toaster } from '../components/core';

const SsoCallback = ({
  route: {
    data: { params, context },
  },
}) => {
  const navigation = useNavigation();
  const [firstRender, setFirstRender] = useState(true);
  const pathComponents = window.location.pathname
    .split('/')
    .filter(x => Boolean(x));
  const pathComponentsUntilSSO = pathComponents.slice(
    0,
    pathComponents.length - 1
  );

  const submit = async data => {
    try {
      switch (pathComponentsUntilSSO[pathComponentsUntilSSO.length - 1]) {
        case 'login':
          await api.loginSSO(data);
          break;
        case 'signup':
          await api.signupSSO(data);
          break;
        default:
          toaster.danger('Internal error, invalid redirect type');
          navigation.navigate('..');
          return;
      }

      const response = await api.user();
      context.setCurrentUser(response.data);
      navigation.navigate(
        params.redirectTo ? decodeURIComponent(params.redirectTo) : '/projects'
      );
    } catch (error) {
      var respCode = error.response.status;
      var respMessage = error.response.data;
      if (respCode == 404) {
        respMessage = 'User not found. Have you already signed up?';
      }

      toaster.danger(`Error code ${respCode}: ${respMessage}`);
      navigation.navigate('..');
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
      toaster.danger(`SSO data not returned.`);
      navigation.navigate('..');
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
      <Text color={'white'} fontSize={6} fontWeight={2} marginLeft={2}>
        Loading...
      </Text>

      <Row marginTop={8}>
        <Button
          justifyContent="center"
          title="Go back"
          onClick={() => {
            navigation.navigate('..');
          }}
        />
      </Row>
    </Column>
  );
};

export default SsoCallback;
