import React from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { useRequest, endpoints } from '../../api';
import Card from '../../components/card';
import Field from '../../components/field';
import Alert from '../../components/alert';
import { Form } from '../../components/core';

const validationSchema = yup.object().shape({
  rawConfig: yup.string().required(),
});
const errorMessages = {
  default: 'Release creation failed.',
};

const CreateRelease = ({
  route: {
    data: { params },
  },
}) => {
  const { data: application } = useRequest(
    endpoints.application({
      projectId: params.project,
      applicationId: params.application,
    }),
    {
      suspense: true,
    }
  );
  const navigation = useNavigation();

  return (
    <Card title="Create Release" size="xlarge">
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        endpoint={endpoints.createRelease({
          projectId: params.project,
          applicationId: application.id,
        })}
        onSuccess={() => {
          navigation.navigate(
            `/${params.project}/applications/${application.name}`
          );
        }}
        onCancel={`/${params.project}/applications/${application.name}/releases`}
        validationSchema={validationSchema}
        defaultValues={{
          rawConfig: application.latestRelease
            ? application.latestRelease.rawConfig
            : '',
        }}
        errorMessages={errorMessages}
        submitLabel="Create"
      >
        <Field type="editor" label="Config" width="100%" name="rawConfig" />
      </Form>
    </Card>
  );
};

export default CreateRelease;
