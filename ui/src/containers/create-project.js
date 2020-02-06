import React, { useEffect } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { endpoints, useMutation } from '../api';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import { Form } from '../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
});

const ProjectCreate = () => {
  const [createProject, { data: project, error }] = useMutation(
    endpoints.createProject(),
    {
      errors: { default: 'Project creation failed.' },
      track: 'Project Created',
    }
  );
  const navigation = useNavigation();

  useEffect(() => {
    if (project) {
      navigation.navigate(`/${project.name}`);
    }
  }, [project]);

  return (
    <Layout alignItems="center">
      <Card title="Create Project" size="medium">
        <Form
          onSubmit={createProject}
          error={error}
          validationSchema={validationSchema}
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
