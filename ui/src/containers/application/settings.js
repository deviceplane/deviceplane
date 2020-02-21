import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { useRequest, endpoints } from '../../api';
import validators from '../../validators';
import Card from '../../components/card';
import Popup from '../../components/popup';
import Field from '../../components/field';
import { Text, Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
});
const updateErrorMessages = {
  default: 'Application update failed.',
};
const deleteErrorMessages = {
  default: 'Application deletion failed.',
};

const ApplicationSettings = ({
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
  const [showDeletePopup, setShowDeletePopup] = useState();
  const navigation = useNavigation();

  return (
    <>
      <Card
        width="100%"
        maxWidth="575px"
        title="Application Settings"
        actions={[
          {
            title: 'Delete',
            onClick: () => setShowDeletePopup(true),
            variant: 'danger',
          },
        ]}
      >
        <Form
          endpoint={endpoints.updateApplication({
            projectId: params.project,
            applicationId: application.id,
          })}
          onSuccess={({ name }) => {
            toaster.success('Application updated.');
            navigation.navigate(`/${params.project}/applications/${name}`);
          }}
          defaultValues={{
            name: application.name,
            description: application.description,
          }}
          validationSchema={validationSchema}
          errorMessages={updateErrorMessages}
        >
          <Field required label="Name" name="name" />
          <Field type="textarea" label="Description" name="description" />
        </Form>
      </Card>
      <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
        <Card title="Delete Application" border size="large">
          <Form
            endpoint={endpoints.deleteApplication({
              projectId: params.project,
              applicationId: application.name,
            })}
            onSuccess={() => {
              setShowDeletePopup(false);
              navigation.navigate(`/${params.project}/applications`);
              toaster.success('Application deleted.');
            }}
            errorMessages={deleteErrorMessages}
            submitLabel="Delete"
            submitVariant="danger"
          >
            <Text>
              You are about to delete the <strong>{application.name}</strong>{' '}
              application.
            </Text>
          </Form>
        </Card>
      </Popup>
    </>
  );
};

export default ApplicationSettings;
