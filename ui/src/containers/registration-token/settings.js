import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';

import { useRequest, endpoints } from '../../api';
import validators from '../../validators';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import { Text, Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  maxRegistrations: yup
    .mixed()
    .notRequired()
    .nullable(),
});
const updateErrorMessages = {
  default: 'Registration token update failed.',
};
const deleteErrorMessages = {
  default: 'Registration token deletion failed.',
};

const RegistrationTokenSettings = ({
  route: {
    data: { params },
  },
}) => {
  const { data: registrationToken } = useRequest(
    endpoints.registrationToken({
      projectId: params.project,
      tokenId: params.token,
    }),
    {
      suspense: true,
    }
  );

  const navigation = useNavigation();
  const [showDeletePopup, setShowDeletePopup] = useState();

  return (
    <Card
      title="Registration Token Settings"
      maxWidth="560px"
      actions={[
        {
          title: 'Delete',
          onClick: () => setShowDeletePopup(true),
          variant: 'danger',
        },
      ]}
    >
      <Form
        endpoint={endpoints.updateRegistrationToken({
          projectId: params.project,
          tokenId: registrationToken.id,
        })}
        onSuccess={() => {
          navigation.navigate(`/${params.project}/provisioning`);
          toaster.success('Registration token updated.');
        }}
        onData={data => ({
          ...data,
          settings: registrationToken.settings,
          maxRegistrations: Number.parseInt(data.maxRegistrations),
        })}
        defaultValues={{
          name: registrationToken.name,
          description: registrationToken.description,
          maxRegistrations: registrationToken.maxRegistrations,
        }}
        errorMessages={updateErrorMessages}
        validationSchema={validationSchema}
      >
        <Field label="Name" name="name" />
        <Field type="textarea" label="Description" name="description" />
        <Field
          type="number"
          label="Maximum Device Registrations"
          name="maxRegistrations"
          description="Limits the number of devices that can be registered using this token."
          hint="Leave empty to allow unlimited registrations."
        />
      </Form>
      <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
        <Card title="Delete Registration Token" border size="large">
          <Form
            endpoint={endpoints.deleteRegistrationToken({
              projectId: params.project,
              tokenId: registrationToken.id,
            })}
            onSuccess={() => {
              setShowDeletePopup(false);
              toaster.success('Registration token deleted.');
              navigation.navigate(`/${params.project}/provisioning`);
            }}
            deleteErrorMessages={deleteErrorMessages}
            submitLabel="Delete"
            submitVariant="danger"
          >
            <Text marginBottom={6}>
              You are about to delete the{' '}
              <strong>{registrationToken.name}</strong> registration token.
            </Text>
          </Form>
        </Card>
      </Popup>
    </Card>
  );
};

export default RegistrationTokenSettings;
