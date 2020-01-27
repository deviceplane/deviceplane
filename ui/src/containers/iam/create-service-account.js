import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import useForm from 'react-hook-form';
import { toaster } from 'evergreen-ui';
import * as yup from 'yup';

import utils from '../../utils';
import api from '../../api';
import validators from '../../validators';
import Field from '../../components/field';
import Card from '../../components/card';
import Alert from '../../components/alert';
import { Row, Form, Button } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
});

const CreateServiceAccount = ({
  route: {
    data: { params },
  },
}) => {
  const { register, handleSubmit, errors } = useForm({ validationSchema });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.createServiceAccount({ projectId: params.project, data });
      toaster.success('Service acccount created successfully.');
      navigation.navigate(`/${params.project}/iam/service-accounts/`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Service Account was not created.');
      console.log(error);
    }
  };

  return (
    <Card title="Create Service Account" size="large">
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
          type="textarea"
          label="Description"
          name="description"
          ref={register}
          errors={errors.description}
        />
        <Button marginTop={3} title="Create" type="submit" />
      </Form>
      <Row marginTop={4}>
        <Button
          title="Cancel"
          variant="text"
          href={`/${params.project}/iam/service-accounts/`}
        />
      </Row>
    </Card>
  );
};

export default CreateServiceAccount;
