import React, { useState } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import { toaster } from 'evergreen-ui';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Row, Button, Form } from '../components/core';

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
  const { register, handleSubmit, errors } = useForm({
    validationSchema,
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.createRegistrationToken({
        projectId: params.project,
        data,
      });
      toaster.success('Registration Token created successfully.');
      navigation.navigate(`/${params.project}/provisioning`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Registration Token was not updated.');
      console.log(error);
    }
  };

  return (
    <Layout alignItems="center">
      <Card title="Create Registration Token" size="large">
        <Alert show={backendError} variant="error" description={backendError} />
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Field
            required
            autoFocus
            label="Name"
            name="name"
            ref={register}
            errors={errors.name}
          />
          <Field
            label="Description"
            type="textarea"
            name="description"
            ref={register}
            errors={errors.description}
          />
          <Field
            type="number"
            label="Maximum Device Registrations"
            name="maxRegistrations"
            description="Limits the number of devices that can be registered using this token."
            hint="Leave empty to allow unlimited registrations."
            errors={errors.maxRegistrations}
            ref={register}
          />
          <Button marginTop={3} title="Create" type="submit" />
        </Form>

        <Row marginTop={4}>
          <Button
            title="Cancel"
            variant="text"
            href={`/${params.project}/provisioning`}
          />
        </Row>
      </Card>
    </Layout>
  );
};

export default CreateRegistrationToken;
