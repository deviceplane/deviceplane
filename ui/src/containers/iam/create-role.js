import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints } from '../../api';
import segment from '../../lib/segment';
import validators from '../../validators';
import Card from '../../components/card';
import Field from '../../components/field';
import { Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  config: yup.string().required(),
});
const errorMessages = {
  default: 'Role creation failed.',
};

const CreateRole = ({
  route: {
    data: { params },
  },
}) => {
  const navigation = useNavigation();
  return (
    <Card title="Create Role" size="large">
      <Form
        endpoint={endpoints.createRole({ projectId: params.project })}
        onSuccess={() => {
          segment.track('Role Created');
          navigation.navigate(`/${params.project}/iam/roles`);
          toaster.success('Role created.');
        }}
        onCancel={`/${params.project}/iam/roles`}
        validationSchema={validationSchema}
        errorMessages={errorMessages}
        submitLabel="Create"
      >
        <Field required autoFocus label="Name" name="name" />
        <Field type="textarea" label="Description" name="description" />
        <Field type="editor" label="Config" name="config" width="100%" />
      </Form>
    </Card>
  );
};

export default CreateRole;
