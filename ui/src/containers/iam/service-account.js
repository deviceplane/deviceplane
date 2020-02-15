import React, { useState, useMemo, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigation } from 'react-navi';
import * as yup from 'yup';
import moment from 'moment';

import api from '../../api';
import utils from '../../utils';
import validators from '../../validators';
import Card from '../../components/card';
import Field from '../../components/field';
import Popup from '../../components/popup';
import Table from '../../components/table';
import Alert from '../../components/alert';
import {
  Row,
  Text,
  Button,
  Form,
  Label,
  Code,
  Icon,
  toaster,
} from '../../components/core';

const validationSchema = yup.object().shape({
  name: validators.name.required(),
  description: yup.string(),
  roles: yup.object(),
});

const ServiceAccount = ({
  route: {
    data: { params, serviceAccount, roles },
  },
}) => {
  const { register, handleSubmit, errors, formState, control } = useForm({
    validationSchema,
    defaultValues: {
      name: serviceAccount.name,
      description: serviceAccount.description,
      roles: roles.reduce(
        (obj, role) => ({
          ...obj,
          [role.name]: !!serviceAccount.roles.find(
            ({ name }) => name === role.name
          ),
        }),
        {}
      ),
    },
  });
  const navigation = useNavigation();
  const [showDeletePopup, setShowDeletePopup] = useState();
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    let error = false;

    try {
      await api.updateServiceAccount({
        projectId: params.project,
        serviceId: serviceAccount.id,
        data,
      });
    } catch (e) {
      error = true;
      console.error(e);
    }

    const roleArray = Object.keys(data.roles);
    for (let i = 0; i < roleArray.length; i++) {
      const role = roleArray[i];
      const choseRole = data.roles[role];
      const hasRole = serviceAccount.roles.find(({ name }) => name === role);
      const roleId = roles.find(({ name }) => name === role).id;
      if (choseRole && !hasRole) {
        try {
          await api.addServiceAccountRoleBindings({
            projectId: params.project,
            serviceId: serviceAccount.id,
            roleId,
          });
        } catch (e) {
          error = true;
          console.error(e);
        }
      } else if (!choseRole & hasRole) {
        try {
          await api.removeServiceAccountRoleBindings({
            projectId: params.project,
            serviceId: serviceAccount.id,
            roleId,
          });
        } catch (e) {
          error = true;
          console.error(e);
        }
      }
    }

    if (error) {
      toaster.danger('Service account update failed.');
    } else {
      toaster.success('Service account updated.');
      navigation.navigate(`/${params.project}/iam/service-accounts`);
    }
  };

  const submitDelete = async () => {
    setBackendError(null);
    try {
      await api.deleteServiceAccount({
        projectId: params.project,
        serviceId: serviceAccount.id,
      });
      toaster.success('Service account deleted.');
      navigation.navigate(`/${params.project}/iam/service-accounts`);
    } catch (error) {
      setBackendError(
        utils.parseError(error, 'Service account deletion failed.')
      );
      console.error(error);
    }
    setShowDeletePopup(false);
  };

  return (
    <>
      <Card
        title={serviceAccount.name}
        subtitle={serviceAccount.description}
        size="xlarge"
        actions={[
          {
            title: 'Delete',
            onClick: () => setShowDeletePopup(true),
            variant: 'danger',
          },
        ]}
        marginBottom={6}
      >
        <Alert variant="error" show={backendError} description={backendError} />
        <Form
          onSubmit={e => {
            setBackendError(null);
            handleSubmit(submit)(e);
          }}
        >
          <Field label="Name" name="name" ref={register} errors={errors.name} />
          <Field
            type="textarea"
            label="Description"
            name="description"
            ref={register}
            errors={errors.description}
          />
          <Label>Choose Individual Roles</Label>
          {roles.map(role => (
            <Field
              multi
              type="checkbox"
              key={role.id}
              label={role.name}
              name={`roles[${role.name}]`}
              control={control}
              errors={errors.roles && errors.roles[role.name]}
            />
          ))}
          <Button
            marginTop={3}
            title="Update"
            type="submit"
            disabled={!formState.dirty}
          />
        </Form>
      </Card>

      <ServiceAccountAccessKeys
        projectId={params.project}
        serviceAccount={serviceAccount}
      />

      <Popup show={showDeletePopup} onClose={() => setShowDeletePopup(false)}>
        <Card title="Delete Service Account" border size="large">
          <Text>
            You are about to delete the <strong>{serviceAccount.name}</strong>{' '}
            service account.
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

export default ServiceAccount;

const ServiceAccountAccessKeys = ({ projectId, serviceAccount }) => {
  const [accessKeys, setAccessKeys] = useState([]);
  const [newAccessKey, setNewAccessKey] = useState();
  const [backendError, setBackendError] = useState();
  const [keyToDelete, setKeyToDelete] = useState();

  const columns = useMemo(
    () => [
      { Header: 'Access Key ID', accessor: 'id' },
      {
        Header: 'Created At',
        accessor: ({ createdAt }) =>
          createdAt ? moment(createdAt).fromNow() : '-',
      },
      {
        Header: ' ',
        Cell: ({ row }) =>
          keyToDelete === row.original.id ? (
            <>
              <Button
                title={<Icon icon="tick-circle" size={16} color="primary" />}
                variant="icon"
                marginRight={4}
                onClick={() => deleteAccessKey(keyToDelete)}
              />
              <Button
                title={<Icon icon="cross" size={16} color="white" />}
                variant="icon"
                onClick={() => setKeyToDelete(null)}
              />
            </>
          ) : (
            <Button
              title={<Icon icon="trash" size={16} color="red" />}
              variant="icon"
              onClick={() => setKeyToDelete(row.original.id)}
            />
          ),
        cellStyle: {
          justifyContent: 'flex-end',
        },
      },
    ],
    [keyToDelete]
  );
  const tableData = useMemo(() => accessKeys, [accessKeys]);

  const fetchAccessKeys = async () => {
    try {
      const response = await api.serviceAccountAccessKeys({
        projectId,
        serviceId: serviceAccount.id,
      });
      setAccessKeys(response.data);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchAccessKeys();
  }, []);

  const createAccessKey = async () => {
    setBackendError(null);
    try {
      const response = await api.createServiceAccountAccessKey({
        projectId,
        serviceId: serviceAccount.id,
      });
      setAccessKeys([response.data, ...accessKeys]);
      setNewAccessKey(response.data.value);
      toaster.success('Access key created.');
    } catch (error) {
      setBackendError(utils.parseError(error, 'Access key creation failed.'));
      console.error(error);
    }
  };

  const deleteAccessKey = async id => {
    setBackendError(null);
    try {
      await api.deleteServiceAccountAccessKey({
        projectId,
        serviceId: serviceAccount.id,
        accessKeyId: id,
      });
      toaster.success('Access key deleted.');
      await fetchAccessKeys();
      setKeyToDelete(null);
    } catch (error) {
      setBackendError(utils.parseError(error, 'Access key deletion failed.'));
      console.error(error);
      setKeyToDelete(null);
    }
  };

  return (
    <Card
      title="Access Keys"
      size="xlarge"
      actions={[{ title: 'Create Access Key', onClick: createAccessKey }]}
    >
      <Alert show={backendError} variant="error" description={backendError} />
      <Alert
        show={!!newAccessKey}
        title="Access Key Created"
        description="Save this key! This is the only time you'll be able to view it. If you lose it, you'll need to create a new access key."
      >
        <Label>Access Key</Label>
        <Row>
          <Code>{newAccessKey}</Code>
        </Row>
      </Alert>
      <Table
        columns={columns}
        data={tableData}
        placeholder={
          <Text>
            There are no <strong>Access Keys</strong>.
          </Text>
        }
      />
    </Card>
  );
};
