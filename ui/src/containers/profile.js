import React from 'react';
import * as yup from 'yup';
import { useCurrentRoute } from 'react-navi';

import { endpoints } from '../api';
import Card from '../components/card';
import Field from '../components/field';
import { Form, toaster } from '../components/core';

const endpoint = endpoints.updateUser();
const validationSchema = yup.object().shape({
  fullName: yup
    .string()
    .required()
    .max(128),
});

const Profile = ({ close }) => {
  const {
    data: {
      context: { currentUser, setCurrentUser },
    },
  } = useCurrentRoute();

  return (
    <Card title="Profile" border size="large">
      <Form
        endpoint={endpoint}
        onSuccess={user => {
          setCurrentUser({ ...currentUser, ...user });
          toaster.success('Profile updated.');
          close();
        }}
        onData={({ fullName }) => {
          let firstName, lastName;
          const firstSpace = fullName.indexOf(' ');
          if (firstSpace === -1) {
            firstName = fullName;
            lastName = ' ';
          } else {
            firstName = fullName.substr(0, firstSpace);
            lastName = fullName.substr(firstSpace + 1);
          }
          return {
            firstName,
            lastName,
          };
        }}
        validationSchema={validationSchema}
        defaultValues={{
          fullName: `${currentUser.firstName} ${currentUser.lastName}`,
        }}
        errorMessages={{ default: 'Profile update failed.' }}
      >
        <Field
          required
          autoCapitalize="words"
          label="Full Name"
          name="fullName"
        />
      </Form>
    </Card>
  );
};

export default Profile;
