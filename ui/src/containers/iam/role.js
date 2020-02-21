import React, { useState } from 'react';
import * as yup from 'yup';
import { useNavigation } from 'react-navi';

import { useRequest, endpoints } from '../../api';
import validators from '../../validators';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import { Text, Form, toaster } from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  config: yup.string().required(),
});
const updateErrorMessages = { default: 'Role update failed.' };
const deleteErrorMessages = { default: 'Role deletion failed.' };

const Role = ({
  route: {
    data: { params },
  },
}) => {
  const { data: role } = useRequest(
    endpoints.role({ projectId: params.project, roleId: params.role }),
    { suspense: true }
  );
  const navigation = useNavigation();
  const [showDeletePopup, setShowDeletePopup] = useState();

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
        <Form
          endpoint={endpoints.updateRole({
            projectId: params.project,
            roleId: role.id,
          })}
          onSuccess={() => {
            toaster.success('Role updated.');
            navigation.navigate(`/${params.project}/iam/roles`);
          }}
          validationSchema={validationSchema}
          defaultValues={{
            name: role.name,
            description: role.description,
            config: role.config,
          }}
          errorMessages={updateErrorMessages}
        >
          <Field required label="Name" name="name" />
          <Field type="textarea" label="Description" name="description" />
          <Field type="editor" label="Config" name="config" width="100%" />
        </Form>
      </Card>
      <Popup
        show={showDeletePopup}
        title="Delete Role"
        onClose={() => setShowDeletePopup(false)}
      >
        <Card title="Delete Role" border size="large">
          <Form
            endpoint={endpoints.deleteRole({
              projectId: params.project,
              roleId: role.id,
            })}
            onSuccess={() => {
              setShowDeletePopup(false);

              navigation.navigate(`/${params.project}/iam/roles`);
              toaster.success('Role deleted.');
            }}
            submitLabel="Delete"
            submitVariant="danger"
            errorMessages={deleteErrorMessages}
          >
            <Text>
              You are about to delete the <strong>{role.name}</strong> role.
            </Text>
          </Form>
        </Card>
      </Popup>
    </>
  );
};

export default Role;
