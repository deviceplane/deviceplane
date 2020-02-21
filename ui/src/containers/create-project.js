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

const endpoint = endpoints.createProject();
const validationSchema = yup.object().shape({
  name: validators.name.required(),
});

const ProjectCreate = () => {
  const navigation = useNavigation();

  return (
    <Layout alignItems="center">
      <Card title="Create Project" size="medium">
        <Form
          endpoint={endpoint}
          onSuccess={project => {
            navigation.navigate(`/${project.name}`);
            segment.track('Project Created');
          }}
          validationSchema={validationSchema}
          errorMessages={{ default: 'Project creation failed.' }}
          submitLabel="Create"
          onCancel="/projects"
        >
          <Field required autoFocus label="Name" name="name" />
        </Form>
      </Card>
    </Layout>
  );
};

export default ProjectCreate;
