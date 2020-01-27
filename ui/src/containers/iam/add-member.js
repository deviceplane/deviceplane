import React, { useState } from 'react';
import { useNavigation } from 'react-navi';
import useForm from 'react-hook-form';

import { toaster } from 'evergreen-ui';

import api from '../../api';
import utils from '../../utils';
import Card from '../../components/card';
import Field from '../../components/field';
import Alert from '../../components/alert';
import { Label, Row, Form, Button, Checkbox } from '../../components/core';

const AddMember = ({
  route: {
    data: { params, roles },
  },
}) => {
  const navigation = useNavigation();
  const { register, handleSubmit, setValue } = useForm({
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
            console.log(e);
          }
        }
      }

      if (error) {
        toaster.warning(
          'Member was added successfully, but roles for the member were not updated properly. Please check the roles of the member.'
        );
      } else {
        navigation.navigate(`/${params.project}/iam/members`);
        toaster.success('Member was added successfully.');
      }
    } catch (error) {
      setBackendError(utils.parseError(error));
      toaster.danger('Member was not added.');
      console.log(error);
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
            key={role.id}
            name={`roles[${role.name}]`}
            as={<Checkbox label={role.name} />}
            register={register}
            setValue={setValue}
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
