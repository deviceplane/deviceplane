import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import useForm from 'react-hook-form';
import { toaster } from 'evergreen-ui';

import api from '../api';
import utils from '../utils';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import Alert from '../components/alert';
import { Button, Row, Form } from '../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
});

const ProjectCreate = () => {
  const navigation = useNavigation();
  const { register, handleSubmit, errors } = useForm({
    validationSchema,
  });
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.createProject(data);
      navigation.navigate(`/${data.name}`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Project was not created.');
      console.log(error);
    }
  };

  return (
    <Layout alignItems="center">
      <Card title="Create Project">
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
          <Button marginTop={3} type="submit" title="Create" />
          <Row marginTop={4}>
            <Button title="Cancel" variant="text" href="/projects" />
          </Row>
        </Form>
      </Card>
    </Layout>
  );
};

export default ProjectCreate;
