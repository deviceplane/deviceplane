import React, { useState } from 'react';
import useForm from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import { toaster } from 'evergreen-ui';

import api from '../../api';
import utils from '../../utils';
import Card from '../../components/card';
import Editor from '../../components/editor';
import Field from '../../components/field';
import Alert from '../../components/alert';
import { Form, Row, Button } from '../../components/core';

const validationSchema = yup.object().shape({
  rawConfig: yup.string().required(),
});

const CreateRelease = ({
  route: {
    data: { params, application },
  },
}) => {
  const { register, handleSubmit, setValue, errors } = useForm({
    validationSchema,
    defaultValues: {
      rawConfig: application.latestRelease
        ? application.latestRelease.rawConfig
        : '',
    },
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      await api.createRelease({
        projectId: params.project,
        applicationId: application.id,
        data,
      });
      navigation.navigate(
        `/${params.project}/applications/${application.name}`
      );
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Release was not created.');
      console.log(error);
    }
  };

  return (
    <Card title="Create Release" size="xlarge">
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        onSubmit={e => {
          setBackendError(null);
          handleSubmit(submit)(e);
        }}
      >
        <Field
          as={<Editor width="100%" />}
          label="Config"
          name="rawConfig"
          register={register}
          setValue={setValue}
          errors={errors.rawConfig}
        />
        <Button type="submit" title="Create" />
      </Form>
      <Row marginTop={4}>
        <Button
          title="Cancel"
          variant="text"
          href={`/${params.project}/applications/${application.name}/releases`}
        />
      </Row>
    </Card>
  );
};

export default CreateRelease;
