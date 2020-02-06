import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints } from '../api';
import segment from '../lib/segment';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import { Form, toaster } from '../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  maxRegistrations: yup
    .mixed()
    .notRequired()
    .nullable(),
});

const CreateRegistrationToken = ({
  route: {
    data: { params },
  },
}) => {
  const navigation = useNavigation();
  return (
    <Layout alignItems="center">
      <Card title="Create Registration Token" size="large">
        <Form
          endpoint={endpoints.createRegistrationToken({
            projectId: params.project,
          })}
          onData={data => ({
            ...data,
            maxRegistrations: Number.parseInt(data.maxRegistrations),
          })}
          onSuccess={() => {
            navigation.navigate(`/${params.project}/provisioning`);
            toaster.success('Registration token created.');
            segment.track('Registration Token Created');
          }}
          onCancel={`/${params.project}/provisioning`}
          errorMessages={{ default: 'Registration token creation failed.' }}
          submitLabel="Create"
          validationSchema={validationSchema}
        >
          <Field required autoFocus label="Name" name="name" />
          <Field label="Description" type="textarea" name="description" />
          <Field
            type="number"
            label="Maximum Device Registrations"
            name="maxRegistrations"
            description="Limits the number of devices that can be registered using this token."
            hint="Leave empty to allow unlimited registrations."
          />
        </Form>
      </Card>
    </Layout>
  );
};

export default CreateRegistrationToken;
