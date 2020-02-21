import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import api, { useRequest, endpoints } from '../api';
import storage from '../storage';
import validators from '../validators';
import Layout from '../components/layout';
import Card from '../components/card';
import Field from '../components/field';
import Popup from '../components/popup';
import { Text, Form, toaster } from '../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
});

const ProjectSettings = ({
  route: {
    data: { params },
  },
}) => {
  const { data: project } = useRequest(
    endpoints.project({ projectId: params.project }),
    {
      suspense: true,
    }
  );
  const navigation = useNavigation();
  const [showDeletePopup, setShowDeletePopup] = useState();

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
              onClick: () => setShowDeletePopup(true),
            },
          ]}
        >
          <Form
            endpoint={endpoints.updateProject({ projectId: project.name })}
            onData={data => ({ ...data, datadogApiKey: project.datadogApiKey })}
            onSuccess={data => {
              console.log(data);
              storage.set('enableSSHKeys', data.enableSSHKeys, data.name);
              navigation.navigate(`/${data.name}`);
              toaster.success('Project updated.');
            }}
            validationSchema={validationSchema}
            defaultValues={{
              name: project.name,
              enableSSHKeys:
                storage.get('enableSSHKeys', project.name) || false,
            }}
            errorMessages={{ default: 'Project update failed.' }}
          >
            <Field required label="Name" name="name" />
            <Field
              type="checkbox"
              label="Enable SSH Keys"
              name="enableSSHKeys"
            />
          </Form>
        </Card>

        <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
          <Card title="Delete Project" border size="large">
            <Text marginBottom={4}>
              This action <strong>cannot</strong> be undone. This will
              permanently delete the <strong>{project.name}</strong> project.
              <p></p>Please enter the project name to confirm.
            </Text>
            <Form
              endpoint={endpoints.deleteProject({ projectId: project.name })}
              onSuccess={() => {
                setShowDeletePopup(false);
                navigation.navigate(`/projects`);
                toaster.success('Project deleted.');
              }}
              submitLabel="Delete"
              errorMessages={{ default: 'Project deletion failed.' }}
            >
              <Field name="name" label="Project Name" />
            </Form>
          </Card>
        </Popup>
      </>
    </Layout>
  );
};

export default ProjectSettings;
