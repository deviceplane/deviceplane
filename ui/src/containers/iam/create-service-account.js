import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints } from '../../api';
import segment from '../../lib/segment';
import validators from '../../validators';
import Field from '../../components/field';
import Card from '../../components/card';
import { Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
});

const CreateServiceAccount = ({
  route: {
    data: { params },
  },
}) => {
  const navigation = useNavigation();

  return (
    <Card title="Create Service Account" size="large">
      <Form
        endpoint={endpoints.createServiceAccount({ projectId: params.project })}
        onSuccess={() => {
          segment.track('Service Account Created');
          navigation.navigate(`/${params.project}/iam/service-accounts/`);
          toaster.success('Service acccount created.');
        }}
        onCancel={`/${params.project}/iam/service-accounts/`}
        validationSchema={validationSchema}
        errorMessages={{ default: 'Service Account creation failed.' }}
      >
        <Field required autoFocus label="Name" name="name" />
        <Field type="textarea" label="Description" name="description" />
      </Form>
    </Card>
  );
};

export default CreateServiceAccount;
