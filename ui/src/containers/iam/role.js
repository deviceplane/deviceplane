import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import * as yup from 'yup';
import { useNavigation } from 'react-navi';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import useToggle from '../../hooks/useToggle';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import Alert from '../../components/alert';
import { Text, Button, Form, toaster } from '../../components/core';

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
  const { register, handleSubmit, errors, formState, control } = useForm({
    validationSchema,
    defaultValues: {
      name: role.name,
      description: role.description,
      config: role.config,
    },
  });
  const navigation = useNavigation();
  const [backendError, setBackendError] = useState();
  const [isDeletePopup, toggleDeletePopup] = useToggle();

  const submit = async data => {
    try {
      await api.updateRole({
        projectId: params.project,
        roleId: role.id,
        data,
      });
      toaster.success('Role updated.');
      navigation.navigate(`/${params.project}/iam/roles`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Role update failed.'));
      console.error(error);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteRole({ projectId: params.project, roleId: role.id });
      toaster.success('Role deleted.');
      navigation.navigate(`/${params.project}/iam/roles`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Role deletion failed.'));
      console.error(error);
    }
    toggleDeletePopup();
  };

  return (
    <>
      <Card
        title={role.name}
        size="large"
        actions={[
          {
            title: 'Delete',
            onClick: toggleDeletePopup,
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
            type="editor"
            label="Config"
            name="config"
            width="100%"
            control={control}
            errors={errors.config}
          />
          <Button
            marginTop={3}
            title="Update"
            type="submit"
            disabled={!formState.dirty}
          />
        </Form>
      </Card>
      <Popup
        show={isDeletePopup}
        title="Delete Role"
        onClose={toggleDeletePopup}
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
