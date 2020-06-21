import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import api from '../api';
import utils from '../utils';
import storage from '../storage';
import validators from '../validators';
import useToggle from '../hooks/useToggle';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import Popup from '../components/popup';
import Alert from '../components/alert';
import {
  Text,
  Button,
  Form,
  Input,
  Label,
  Group,
  toaster,
} from '../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
});

const ProjectSettings = ({
  route: {
    data: { params, project },
  },
}) => {
  const { control, register, handleSubmit, errors, formState } = useForm({
    validationSchema,
    defaultValues: {
      name: project.name,
      legacy: storage.get('legacy') || false,
    },
  });
  const navigation = useNavigation();
  const [isDeletePopup, toggleDeletePopup] = useToggle();
  const [confirmation, setConfirmation] = useState();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    data.datadogApiKey = project.datadogApiKey;
    try {
      storage.set('legacy', data.legacy);
      await api.updateProject({ projectId: project.name, data });
      toaster.success('Project updated.');
      navigation.navigate(`/${data.name}`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Project update failed.'));
      console.error(error);
    }
  };

  const submitDelete = async e => {
    e.preventDefault();
    setBackendError(null);
    try {
      await api.deleteProject({ projectId: project.name });
      toaster.success('Project deleted.');
      navigation.navigate(`/projects`);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Project deletion failed.'));
      console.error(error);
    }
    toggleDeletePopup();
  };

  return (
    <Layout alignItems="center">
      <>
        <Card
          title="Project Settings"
          maxWidth="540px"
          actions={[
            {
              title: 'Delete',
              variant: 'danger',
              onClick: toggleDeletePopup,
            },
          ]}
        >
          <Alert
            show={backendError}
            variant="error"
            description={backendError}
          />
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
              type="checkbox"
              label="Enable legacy features"
              name="legacy"
              control={control}
            />
            <Button
              marginTop={3}
              type="submit"
              title="Update"
              disabled={!formState.dirty}
            />
          </Form>
        </Card>

        <Popup show={isDeletePopup} onClose={toggleDeletePopup}>
          <Card title="Delete Project" border size="large">
            <Text marginBottom={4}>
              This action <strong>cannot</strong> be undone. This will
              permanently delete the <strong>{params.project}</strong> project.
              <p></p>Please enter the project name to confirm.
            </Text>
            <Form onSubmit={submitDelete}>
              <Group>
                <Label>Project Name</Label>
                <Input
                  onChange={e => setConfirmation(e.target.value)}
                  value={confirmation}
                />
              </Group>

              <Button
                variant="danger"
                type="submit"
                title="Delete"
                disabled={confirmation !== project.name}
              />
            </Form>
          </Card>
        </Popup>
      </>
    </Layout>
  );
};

export default ProjectSettings;
