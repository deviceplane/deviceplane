import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import { useForm } from 'react-hook-form';

import api from '../../api';
import utils from '../../utils';
import Card from '../../components/card';
import Field from '../../components/field';
import Alert from '../../components/alert';
import { Label, Row, Form, Button, toaster } from '../../components/core';

const AddMember = ({
  route: {
    data: { params, roles },
  },
}) => {
  const navigation = useNavigation();
  const { register, handleSubmit, control } = useForm({
    defaultValues: {
      roles: roles.reduce((obj, role) => ({ ...obj, [role.name]: false }), {}),
    },
  });
  const [backendError, setBackendError] = useState();

  const submit = async data => {
    try {
      const {
        data: { userId },
      } = await api.addMember({
        projectId: params.project,
        data: { email: data.email },
      });

      let error = false;

      const roleArray = Object.keys(data.roles);
      for (let i = 0; i < roleArray.length; i++) {
        const role = roleArray[i];
        const hasRole = data.roles[role];
        if (hasRole) {
          try {
            await api.addMembershipRoleBindings({
              projectId: params.project,
              userId,
              roleId: role,
            });
          } catch (e) {
            error = true;
            console.error(e);
          }
        }
      }

      if (error) {
        setBackendError(
          'Member added, but roles for the member were not updated. Please verify the roles are valid.'
        );
      } else {
        navigation.navigate(`/${params.project}/iam/members`);
        toaster.success('Member added.');
      }
    } catch (error) {
      setBackendError(utils.parseError(error, 'Adding member failed.'));
      console.error(error);
    }
  };

  return (
    <Card title="Add Member" size="large">
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        onSubmit={e => {
          setBackendError(null);
          handleSubmit(submit)(e);
        }}
      >
        <Field
          autoFocus
          required
          label="Email"
          type="email"
          name="email"
          ref={register}
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
          />
        ))}
        <Button marginTop={3} type="submit" title="Add Member" />
      </Form>
      <Row marginTop={4}>
        <Button
          title="Cancel"
          variant="text"
          href={`/${params.project}/iam/members`}
        />
      </Row>
    </Card>
  );
};

export default AddMember;
