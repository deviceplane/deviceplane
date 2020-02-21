import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints } from '../api';
import segment from '../lib/segment';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import { Form } from '../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
});

const CreateApplication = ({
  route: {
    data: { params },
  },
}) => {
  const navigation = useNavigation();

  return (
    <Layout alignItems="center">
      <Card title="Create Application" size="large">
        <Form
          endpoint={endpoints.createApplication({ projectId: params.project })}
          onSuccess={({ name }) => {
            segment.track('Application Created');
            navigation.navigate(`/${params.project}/applications/${name}`);
          }}
          onCancel={() =>
            navigation.navigate(`/${params.project}/applications`)
          }
          validationSchema={validationSchema}
          submitLabel="Create"
          errorMessages={{ default: 'Application creation failed.' }}
        >
          <Field required autoFocus label="Name" name="name" />
          <Field label="Description" name="description" type="textarea" />
        </Form>
      </Card>
    </Layout>
  );
};

export default CreateApplication;
