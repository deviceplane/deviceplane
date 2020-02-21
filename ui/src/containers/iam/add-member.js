import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import { useForm } from 'react-hook-form';

import api, { useRequest, endpoints } from '../../api';
import utils from '../../utils';
import Card from '../../components/card';
import Field from '../../components/field';
import { Form, toaster } from '../../components/core';

const AddMember = ({
  route: {
    data: { params },
  },
}) => {
  const { data: roles } = useRequest(
    endpoints.roles({ projectId: params.project }),
    { suspense: true }
  );

  // const submit = async data => {
  //   try {
  //     const {
  //       data: { userId },
  //     } = await
  //       data: { email: data.email },
  //     });

  //     let error = false;

  //     const roleArray = Object.keys(data.roles);
  //     for (let i = 0; i < roleArray.length; i++) {
  //       const role = roleArray[i];
  //       const hasRole = data.roles[role];
  //       if (hasRole) {
  //         try {
  //           await api.addMembershipRoleBindings({
  //             projectId: params.project,
  //             userId,
  //             roleId: role,
  //           });
  //         } catch (e) {
  //           error = true;
  //           console.error(e);
  //         }
  //       }
  //     }

  // setBackendError(
  //   'Member added, but roles for the member were not updated. Please verify the roles are valid.'
  // );

  return (
    <Card title="Add Member" size="large">
      <Alert show={backendError} variant="error" description={backendError} />
      <Form
        endpoint={endpoints.addMember({
          projectId: params.project,
        })}
        onSuccess={() => {
          segment.track('Member Added');
          navigation.navigate(`/${params.project}/iam/members`);
          toaster.success('Member added.');
        }}
        onCancel={`/${params.project}/iam/members`}
        defaultValues={{
          roles: roles.reduce(
            (obj, role) => ({ ...obj, [role.name]: false }),
            {}
          ),
        }}
        submitLabel="Add Member"
        errorMessages={{ default: 'Adding member failed.' }}
      >
        <Field autoFocus required label="Email" type="email" name="email" />
        <Label>Choose Individual Roles</Label>
        {roles.map(role => (
          <Field
            multi
            type="checkbox"
            key={role.id}
            label={role.name}
            name={`roles[${role.name}]`}
          />
        ))}
      </Form>
    </Card>
  );
};

export default AddMember;
