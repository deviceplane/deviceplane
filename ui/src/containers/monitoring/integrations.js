import React from 'react';
import * as yup from 'yup';
import { useNavigation } from 'react-navi';

import { endpoints } from '../../api';
import { Form, toaster } from '../../components/core';
import Card from '../../components/card';
import Field from '../../components/field';

const validationSchema = yup.object().shape({
  datadogApiKey: yup.string(),
});

const Integrations = ({
  route: {
    data: { params, project },
  },
}) => {
  const navigation = useNavigation();

  return (
    <Card
      title="Integrations"
      size="large"
      subtitle={
        project.datadogApiKey
          ? ''
          : 'A Datadog API key to is required to use monitoring.'
      }
    >
      <Form
        endpoint={endpoints.updateProject({ projectId: params.project })}
        onSuccess={() => {
          if (!project.datadogApiKey) {
            navigation.navigate(`/${params.project}/monitoring/project`);
          } else {
            toaster.success('Integrations updated.');
          }
        }}
        onData={data => ({ ...data, name: params.project })}
        validationSchema={validationSchema}
        defaultValues={{
          datadogApiKey: project.datadogApiKey,
        }}
        errorMessages={{ default: 'Integrations update failed.' }}
      >
        <Field label="Datadog API Key" name="datadogApiKey" />
      </Form>
    </Card>
  );
};

export default Integrations;
