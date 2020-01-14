import React, { useState } from 'react';
import useForm from 'react-hook-form';
import * as yup from 'yup';
import { useNavigation } from 'react-navi';
import { toaster } from 'evergreen-ui';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import Editor from '../../components/editor';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import Alert from '../../components/alert';
import { Text, Button, Form } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  config: yup.string().required(),
});

const Role = ({
  route: {
    data: { params, role },
  },
}) => {
  const { register, handleSubmit, errors, formState, setValue } = useForm({
    validationSchema,
    defaultValues: {
      name: role.name,
      description: role.description,
      config: role.config,
    },
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();
  const [showDeletePopup, setShowDeletePopup] = useState();

  const submit = async data => {
    try {
      await api.updateRole({
        projectId: params.project,
        roleId: role.id,
        data,
      });
      toaster.success('Role was updated successfully.');
      navigation.navigate(`/${params.project}/iam/roles`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Role was not updated.');
      console.log(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteRole({ projectId: params.project, roleId: role.id });
      toaster.success('Successfully deleted role.');
      navigation.navigate(`/${params.project}/iam/roles`);
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Role was not deleted.');
      console.log(error);
    }
    setShowDeletePopup(false);
  };

  return (
    <>
      <Card
        title={role.name}
        size="large"
        actions={[
          {
            title: 'Delete',
            onClick: () => setShowDeletePopup(true),
            variant: 'danger',
          },
        ]}
      >
        <Alert show={backendError} variant="error" description={backendError} />
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Field
            required
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
          <Field
            as={<Editor width="100%" height="160px" />}
            label="Config"
            name="config"
            register={register}
            setValue={setValue}
            errors={errors.config}
          />
          <Button title="Update" type="submit" disabled={!formState.dirty} />
        </Form>
      </Card>
      <Popup
        show={showDeletePopup}
        title="Delete Role"
        onClose={() => setShowDeletePopup(false)}
      >
        <Card title="Delete Role" border size="large">
          <Text>
            You are about to delete the <strong>{role.name}</strong> role.
          </Text>
          <Button
            marginTop={5}
            title="Delete"
            onClick={submitDelete}
            variant="danger"
          />
        </Card>
      </Popup>
    </>
  );
};

export default Role;
